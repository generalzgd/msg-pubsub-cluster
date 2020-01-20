/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @time: 2019/12/6 5:11 下午
 * @project: packagesubscribesvr
 */

package mgr

import (
	`encoding/binary`
	`os`
	`path/filepath`
	`sync`
	`sync/atomic`
	`time`

	`github.com/astaxie/beego/logs`
	`github.com/generalzgd/cluster-plugin/plugin`
	`github.com/generalzgd/comm-libs/number`
	_ `github.com/generalzgd/grpc-svr-frame/grpc-consul`
	ctrl `github.com/generalzgd/grpc-svr-frame/grpc-ctrl`

	`github.com/generalzgd/msg-subscriber/codec/body`
	`github.com/generalzgd/msg-subscriber/codec/head`
	`github.com/generalzgd/msg-subscriber/config`
	`github.com/generalzgd/msg-subscriber/define`
	`github.com/generalzgd/msg-subscriber/iface`
	`github.com/generalzgd/msg-subscriber/iproto`
	`github.com/generalzgd/msg-subscriber/receive`
	`github.com/generalzgd/msg-subscriber/storage`
)

const (
	IndexFlow  = "index"
	ReportFlow = "report_flow"
	LeaderFlow = "leader_flow"
	WorkFlow   = "work_flow"
	DeadFlow   = "dead_flow"
	//
	IndexKey = IndexFlow
)

var (
	cmdDecoder  head.Decoder
	bodyDecoder body.Decoder
)

type Manager struct {
	// 加入集群插件
	*plugin.ClusterPlugin
	ctrl.GrpcController
	cfg                              config.AppConfig
	TcpReceiver                       *receive.TcpReceiver
	imLeader                         bool
	imFollower                       bool
	acceptClientFlag                 int32                         // 是否启动端口侦听，开始接收producer or consumer
	syncClusterSubscribedInfoOnStart int32                         // 集群节点启动的时候，是否同步已订阅信息标记
	CurrentLeaderMeta                plugin.ServerMeta             // 当前集群的leader信息
	rpcHandlers                      map[string]RpcHandler         //
	raftHandlers                     map[iproto.CusCmd]RaftHandler //
	index                            number.UnidirectionalNum      // 全局的消息索引，要求集群内索引保持一致
	indexSyncFlag                    bool                          // 索引同步标记
	indexQueue                       iface.IStoreBridge            // 索引存储
	publishQueue                     iface.IStoreBridge            // 分发流
	deadQueue                        iface.IStoreBridge            // 死信流, todo 是不是应该存文件？永久存储
	allSubscribeMap                  *define.SubscribeFlagMap      // 集群内所有的节点订阅计数
	subscribeMap                     *define.SubscribeMap          // 当前节点的消息订阅映射表
	postId2Subscriber                sync.Map                      // post id 与订阅者的关系
	subscriberTimer                  sync.Map                      // 订阅者的延时清理程序定时器
	askId                            uint32                        //
	quit                             chan struct{}
}

func NewManager() *Manager {
	return &Manager{
		GrpcController:  ctrl.MakeGrpcController(),
		allSubscribeMap: define.NewSubscribeFlagMap(),
		subscribeMap:    define.NewSubscribeMap(),
		quit:            make(chan struct{}),
	}
}

func (p *Manager) Init(cfg config.AppConfig) error {
	p.cfg = cfg
	//
	cmdDecoder = head.NewDecoder(cfg.Decode.CmdPos, cfg.Decode.CmdSize, cfg.Decode.LenPos, cfg.Decode.LenSize, binary.LittleEndian)
	bodyDecoder = body.NewDecoder(cfg.Decode.HeadSize, binary.LittleEndian)
	//
	buckets := []string{IndexFlow, ReportFlow, LeaderFlow, WorkFlow, DeadFlow}
	path := filepath.Join(filepath.Dir(os.Args[0]), "bolt.db")
	storeOperator := storage.NewStoreOperator(cfg.QueueMode, cfg.TreeDegree, path, func(bucket string) iface.StoreItem {
		if bucket == IndexFlow {
			return &define.IndexPack{}
		}
		return define.NewFlowPack(cmdDecoder, bodyDecoder)
	}, buckets...)

	p.indexQueue = storage.NewStoreBridge(IndexFlow, storeOperator)
	// p.leaderQueue = storage.NewStoreBridge(LeaderFlow, storeOperator)
	p.publishQueue = storage.NewStoreBridge(WorkFlow, storeOperator)
	p.deadQueue = storage.NewStoreBridge(DeadFlow, storeOperator)
	//
	if v, _ := p.indexQueue.GetUint64(IndexKey); v > 0 {
		p.index.SetNumber(v)
	}
	//
	/*p.SetGrpcPoolConfig(ctrl.GrpcPoolConfig{
		Init:            cfg.GrpcPool.Init,
		Capacity:        cfg.GrpcPool.Capacity,
		IdleTimeout:     cfg.GrpcPool.IdleTimeout,
		MaxLifeDuration: cfg.GrpcPool.MaxLifeTimeout,
	})*/
	//
	p.initRpcHandlers()
	p.initRaftHandlers()
	//
	cluArgs := plugin.ClusterArgs{
		Role:     cfg.Name,
		Agent:    plugin.NodeType(cfg.Cluster.NodeType),
		SerfPort: cfg.Cluster.SerfPort,
		RaftPort: cfg.Cluster.RaftPort,
		RpcPort:  cfg.Cluster.RpcPort,
		HttpPort: cfg.Cluster.HttpPort,
		Except:   cfg.Cluster.Except,
		Peers:    cfg.Cluster.Peers,
	}
	c, err := plugin.CreateCluster(cluArgs, p, p)
	if err != nil {
		return err
	}
	p.ClusterPlugin = c

	// p.Run()
	return nil
}

func (p *Manager) Destroy() {
	p.DisposeGrpcConn("")
	p.ClusterPlugin.Shutdown()
	p.indexQueue.Close()
	// p.reportQueue.Close()
	p.publishQueue.Close()
	p.deadQueue.Close()
	close(p.quit)
}

// 尝试启动客户端侦听，接收producer and consumer
func (p *Manager) tryStartAccept() {
	if atomic.CompareAndSwapInt32(&p.acceptClientFlag, 0, 1) {
		time.AfterFunc(time.Second*2, func() {
			// 启动post
			if err := p.StartTcpWork(); err != nil {
				logs.Error("tryStartAccept() start post got err=(%v)", err)
				return
			}
			// 启动rpc
			if err := p.StartSubscribeRpcSvr(); err != nil {
				logs.Error("tryStartAccept() start rpc got err=(%v)", err)
				return
			}
		})
	}
}

//
func (p *Manager) Run() {
	// go p.watchLeaderFlow()
	// go p.watchReportFlow()
	go p.watchPublishFlow()
	go p.watchDeadFlow()
}

// 非leader节点处理, leader跳过运行，发消息给订阅者
func (p *Manager) watchPublishFlow() {
	ticker := time.NewTicker(p.cfg.MaxPackDelay)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case <-ticker.C:
			p.pullPublishQueue()
		case <-p.quit:
			return
		}
	}
}

// 考虑所有的节点都有分发功能
func (p *Manager) pullPublishQueue() {
	list, err := p.publishQueue.GetBatch(p.cfg.BatchSize)
	if err != nil {
		logs.Error("pullPublishQueue() got err=%v", err)
		return
	}
	if len(list) == 0 {
		// logs.Debug("pullPublishQueue() got empty")
		return
	}

	p.doPublishMsg(list...)
}

// leader的死信队列是发给其他节点
func (p *Manager) watchDeadFlow() {
	ticker := time.NewTicker(p.cfg.DeadDelay)
	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case <-ticker.C:
			p.pullDeadQueue()
		case <-p.quit:
			return
		}
	}

}

// 考虑所有的节点都有分发功能
func (p *Manager) pullDeadQueue() {
	// if !p.imFollower {
	// 	return
	// }
	list, err := p.deadQueue.GetBatch(p.cfg.DeadBatchSize)
	if err != nil || len(list) == 0 {
		return
	}

	p.doRepublishMsg(list...)
}

// 封装为原子操作, 只接收来自leader的修改
func (p *Manager) SetIndex(v uint64) {
	// 只往大的写
	p.index.SetNumber(v, func(v uint64) {
		p.indexQueue.SetUnit64(IndexKey, v)
	})
}

func (p *Manager) GetIndex() uint64 {
	return p.index.GetNumber() // atomic.LoadUint64(&p.index)
}

func (p *Manager) IncreaseIndex() (uint64, bool) {
	if p.imLeader {
		v := p.index.AutoIncrease(func(v uint64) {
			p.indexQueue.SetUnit64(IndexKey, v)
		})
		return v, true
	}
	return 0, false
}

func (p *Manager) IncreasedAskId() uint32 {
	return atomic.AddUint32(&p.askId, 1)
}

func (p *Manager) CleanSubscriberTimer(key string) bool {
	if v, ok := p.subscriberTimer.Load(key); ok {
		if t, ok := v.(*time.Timer); ok {
			t.Stop()
		}
	}
	return true
}

// func StoreItemToFlowPack(in ...iface.StoreItem) []*define.FlowPack {
// 	out := make([]*define.FlowPack, len(in))
// 	for i, v := range in {
// 		out[i] = v.(*define.FlowPack)
// 	}
// 	return out
// }

func StoreItemToStorePack(in ...iface.StoreItem) []*iproto.StorePack {
	out := make([]*iproto.StorePack, len(in))
	for i, v := range in {
		t := v.(*define.FlowPack)
		p := t.ToStorePack()
		out[i] = &p
	}
	return out
}

func StorePackToStoreItem(in ...*iproto.StorePack) []iface.StoreItem {
	out := make([]iface.StoreItem, len(in))
	for i, v := range in {
		p := define.NewFlowPack(cmdDecoder, bodyDecoder)
		p.FromStorePack(*v)
		out[i] = p
	}
	return out
}

func StoreItemToUint64(in ...iface.StoreItem) []uint64 {
	out := make([]uint64, 0, len(in))
	for _, v := range in {
		out = append(out, v.GetIndex())
	}
	return out
}
