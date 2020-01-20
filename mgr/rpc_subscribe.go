/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: rpc_subscribe.go
 * @time: 2019/12/30 12:33 下午
 * @project: msgsubscribesvr
 */

package mgr

import (
	`context`
	`errors`
	`fmt`
	`net`
	`time`

	`github.com/astaxie/beego/logs`
	`github.com/golang/protobuf/proto`
	`google.golang.org/grpc`
	`google.golang.org/grpc/reflection`

	`github.com/generalzgd/msg-subscriber/codec`
	`github.com/generalzgd/msg-subscriber/define`
	`github.com/generalzgd/msg-subscriber/iproto`
	`github.com/generalzgd/msg-subscriber/sender`
	`github.com/generalzgd/msg-subscriber/util`
)

func (p *Manager) StartSubscribeRpcSvr() error {
	opts := []grpc.ServerOption{
		grpc.MaxConcurrentStreams(1000),
		grpc.MaxRecvMsgSize(32 * 1024),
		grpc.MaxSendMsgSize(32 * 1024),
		grpc.ReadBufferSize(8 * 1024),
		grpc.WriteBufferSize(8 * 1024),
		grpc.ConnectionTimeout(5 * time.Second),
	}

	addr := fmt.Sprintf(":%d", p.cfg.SubscribeCfg.Port)
	s := grpc.NewServer(opts...)
	iproto.RegisterSubscribeSvrServer(s, p)
	reflection.Register(s)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logs.Error("failed to listen: %v", err)
	}
	go func() {
		if err := s.Serve(lis); err != nil {
			logs.Error("failed to serve: %v", err)
		}
	}()
	logs.Debug("start serve subscribe.", addr)
	return nil
}

func (p *Manager) Subscribe(ctx context.Context, req *iproto.SubscribeRequest) (*iproto.SubscribeReply, error) {
	peer, ok := p.GetPeer(ctx)
	if !ok {
		return &iproto.SubscribeReply{}, errors.New("cannot get peer address")
	}
	info := &define.SubscribeInfo{
		FromType:    define.SubscribeTypeRpc,
		ConsumerKey: req.ConsumerKey,
		Sender:      sender.NewGrpcSender(req.Name, peer.Addr.String(), req.RpcAddr, &p.GrpcController),
	}
	ids := util.Uint32ToInt(req.CmdIds...) // make([]uint16, len(req.Cmdid))
	if req.Act == true {
		p.doneSubscribe(ids, info)
	} else {
		p.doneUnsubscribe(ids, info)
	}

	// 测试
	/*time.AfterFunc(time.Second*30, func() {
		bts := []byte(`{"cmdid":"chatmessage"}`)
		pk := gocmd.PostPacket{
			Length:   uint32(len(bts)),
			CmdId:    gocmd.ID_Chatmessage,
			Body:     bts,
		}
		info.Send(&define.FlowPack{
			Index:      1,
			Data:       pk.Serialize(),
			Pack:       pk,
		})
	})*/

	return &iproto.SubscribeReply{}, nil
}

func (p *Manager) Produce(ctx context.Context, req *iproto.ProduceRequest) (rep *iproto.ProduceReply, err error) {
	rep = &iproto.ProduceReply{}

	args := &iproto.ProduceRequest{}
	if err = proto.Unmarshal(req.Data, args); err != nil {
		return
	}
	//
	pack := codec.NewDataPack(args.Data, headDecoder, bodyDecoder)
	cmd, _ := pack.GetHead()
	// 未订阅的消息一律过滤掉
	if p.allSubscribeMap.Has(cmd) {
		var index uint64
		index, err = p.askIncreasedIndex()
		if err != nil {
			logs.Error("doReportMsg() get increased index err. %v", err)
			return
		}

		msg := &define.FlowPack{
			Index: index,
			Pack:  pack,
		}
		p.doReportMsg(msg)
	}
	return
}
