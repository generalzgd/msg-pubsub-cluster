/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: rpc_call.go
 * @time: 2019/12/26 5:36 下午
 * @project: msgsubscribesvr
 */

package mgr

import (
	`sync/atomic`
	`time`

	`github.com/astaxie/beego/logs`
	`github.com/generalzgd/cluster-plugin/plugin`
	`github.com/golang/protobuf/proto`
	`github.com/golang/protobuf/ptypes`

	`github.com/generalzgd/msg-subscriber/define`
	`github.com/generalzgd/msg-subscriber/iface`
	`github.com/generalzgd/msg-subscriber/iproto`
	`github.com/generalzgd/msg-subscriber/util`
)

type ProtoMsgFactory func() proto.Message

// 集群广播访问
func (p *Manager) ClusterRpcCall(cmd string, req proto.Message,
	factory ProtoMsgFactory, stopWhenErr bool, retry int, exceptMe bool) (out []proto.Message, err error) {
	list := p.MemberList()
	out = make([]proto.Message, 0, len(list))
	if retry < 1 {
		retry = 1
	}

	for _, m := range list {
		meta := plugin.ServerMeta{Name: m.Name, Addr: m.Addr}
		if e := meta.FromMap(m.Tags); e != nil {
			err = e
			if stopWhenErr {
				return
			}
			continue
		}
		addr := meta.RpcAddr()
		if exceptMe && addr == p.CurrentServerMeta.RpcAddr() {
			continue
		}
		//logs.Debug("ClusterRpcCall() start rpc call.", addr)
		for i := 0; i < retry; i++ {
			if rep, e := p.RpcCall(addr, cmd, p.IncreasedAskId(), req); e != nil {
				err = e
				if stopWhenErr {
					return
				}
				goto sleep
			} else {
				if rep.Code == 0 && rep.Data != nil {
					rt := factory()
					if e := ptypes.UnmarshalAny(rep.Data, rt); e != nil {
						err = e
						if stopWhenErr {
							return
						}
						goto sleep
					}
					out = append(out, rt)
					continue
				} else {
					err = define.ErrorCode
					goto sleep
				}
			}
		sleep:
			time.Sleep(time.Millisecond * 30)
		}
	}
	return
}

// 收集集群索引
func (p *Manager) collectClusterIndex() (out []*iproto.CollectClusterIndexReply, err error) {
	var tmp []proto.Message
	cmd := iproto.CusCmd_CollectClusterIndex.String()
	retFactory := func() proto.Message {
		return &iproto.CollectClusterIndexReply{}
	}
	req := &iproto.CollectClusterIndexRequest{}

	tmp, err = p.ClusterRpcCall(cmd, req, retFactory, false, p.cfg.Retry, true)
	if err != nil {
		logs.Error("collectClusterIndex() got err=%v", err)
		return
	}

	out = make([]*iproto.CollectClusterIndexReply, 0, len(tmp))
	for _, it := range tmp {
		out = append(out, it.(*iproto.CollectClusterIndexReply))
	}
	return
}

// leader 同步集群索引 Rpc
func (p *Manager) syncClusterIndex(index uint64) (err error) {
	if index == 0 {
		index = p.GetIndex()
	}

	req := &iproto.SyncClusterIndexRequest{
		Index: index,
	}
	if err = p.CallLeaderApply(iproto.CusCmd_SyncClusterIndex, req); err != nil {
		logs.Error("syncClusterIndex() got err=%v", err)
		return
	}
	return
}

// 节点向集群leader获取索引
func (p *Manager) askClusterIndex() (out *iproto.AskIndexReply, err error) {
	if p.imLeader {
		return &iproto.AskIndexReply{
			Res:   &iproto.CommonReply{},
			Index: p.GetIndex(),
		}, nil
	}
	addr := p.CurrentLeaderMeta.RpcAddr()
	req := &iproto.AskIndexRequest{
		Index: p.GetIndex(),
	}
	rep, err := p.RpcCall(addr, iproto.CusCmd_AskIndex.String(), p.IncreasedAskId(), req)
	if err != nil {
		logs.Error("askClusterIndex() got err=%v", err)
		return &iproto.AskIndexReply{}, err
	}

	if rep.Data != nil {
		out = &iproto.AskIndexReply{}
		if err = ptypes.UnmarshalAny(rep.Data, out); err != nil {
			logs.Error("askClusterIndex() got err=%v", err)
			return &iproto.AskIndexReply{}, err
		} else {
			logs.Info("askClusterIndex() got=%v", out)
		}
	}

	return
}

// 请求获取集群自增后的索引
func (p *Manager) askIncreasedIndex() (uint64, error) {
	if p.imLeader {
		v, _ := p.IncreaseIndex()
		if err := p.CallLeaderApply(iproto.CusCmd_SyncClusterIndex, &iproto.SyncClusterIndexRequest{
			Index: v,
		}); err != nil {
			logs.Error("askIncreasedIndex() leader got err=%v", err)
			return 0, err
		}
		return v, nil
	} else if p.imFollower {
		// 请求leader
		cmd := iproto.CusCmd_AskIncreasedIndex.String()
		addr := p.CurrentLeaderMeta.RpcAddr()
		req := &iproto.AskIncreasedIndexRequest{}
		//
		got, err := p.RpcCall(addr, cmd, p.IncreasedAskId(), req)
		if err != nil {
			logs.Error("askIncreasedIndex() follower got err=%v", err)
			return 0, err
		}
		rep := &iproto.AskIncreasedIndexReply{}
		if err := ptypes.UnmarshalAny(got.Data, rep); err != nil {
			logs.Error("askIncreasedIndex() follower got err=%v", err)
			return 0, err
		}
		return rep.Index, nil
	}
	return 0, define.IndexErr
}

// 向leader上报，由leader同步给集群，自己直接调用，默认都是订阅过的消息
func (p *Manager) doReportMsg(list ...iface.StoreItem) error {
	//logs.Debug("reportMsg() list=%v", list)

	if len(list) < 1 {
		return nil
	}

	//
	cmd := iproto.CusCmd_ReportNewMsg
	req := &iproto.NewMsgRequest{
		Data: StoreItemToStorePack(list...),
	}
	if err := p.CallLeaderApply(cmd, req); err != nil {
		logs.Error("doReportMsg() call leader apply err. %v", err)
	}

	return nil
}

// leader 同步集群索引操作
// 1. 广播获取集群每个节点的索引值，
// 2. 如果存在不一致，则取最大值并同步给每个节点
func (p *Manager) doSyncClusterIndex() {
	out, err := p.collectClusterIndex()
	//
	if err != nil {
		logs.Error("cluster: index sync fail with err.", err)
		return
	}
	max := uint64(0)
	same := true
	for i, it := range out {
		if it.Index > max {
			max = it.Index
		}
		if i >= 1 && it.Index != out[0].Index {
			same = false
		}
	}
	myIndex := p.GetIndex()
	// 索引一致
	if max == myIndex {
		return
	}
	// 其他节点的大,
	if max > myIndex {
		// 更新自己的索引
		p.SetIndex(max)
		//
		if !same {
			// 同步其他节点
			p.syncClusterIndex(max)
		}
	} else if max < myIndex {
		// 同步其他节点
		p.syncClusterIndex(myIndex)
	}
}

// follower 向leader获取最新的索引（不自增），
// 1. 广播获取集群每个节点的索引值，
// 2. 如果索引大于本节点，则存起来
func (p *Manager) doAskClusterIndex() {
	res, err := p.askClusterIndex()
	if err != nil {
		logs.Error("doAskClusterIndex() got err=%v", err)
		return
	}

	// 修改本节点索引
	p.SetIndex(res.Index)
}

// 分发消息
// 1. 未在当前节点注册对消息，直接忽略
// 2. 发布的消息，
// 2.1 发送失败，则需要上报给leader, 删除发布队列，进入死信队列，并同步集群
// 2.2 发送成功，上报给leader，删除发布队列，并同步集群
func (p *Manager) doPublishMsg(list ...iface.StoreItem) error {
	if len(list) < 1 {
		return nil
	}
	// 发布成功的数据列表
	okList := make([]iface.StoreItem, 0, len(list))
	// 发布失败等数据列表
	deadList := make([]iface.StoreItem, 0, len(list))
	// 过期未发布的数据列表
	warnList := make([]iface.StoreItem, 0, len(list))
	// 异常数据，要删除列表
	delList := make([]iface.StoreItem, 0, len(list))

	// node := p.CurrentServerMeta.RaftAddr() // 当前节点
	for _, it := range list {
		var infos []iface.IConsumer
		var booked bool
		if id := it.GetCmdId(); id > 0 {
			infos, booked = p.subscribeMap.GetConsumer(id)
		} else {
			logs.Warn("doPublishMsg() cmdid empty")
			continue
		}

		if !booked {
			// 过期未订阅，未发布的数据
			if it.MaxRetry() < 1 && it.HasExpire(p.cfg.MaxPackDelay.Seconds()*3) {
				warnList = append(warnList, it)
			}
			logs.Warn("doPublishMsg() cmdid unbooked=[%v]", it)
			// 再等2轮，等其他节点发布成功，然后删除
			continue
		}

		err := p.publish(infos, it, false)
		if err != nil {
			deadList = append(deadList, it)
			logs.Error("doPublishMsg() got err=[%v]", err)
		} else {
			okList = append(okList, it)
		}
	}

	p.doPublishLiveResult(okList, deadList, warnList, delList)
	return nil
}

func (p *Manager) doPublishLiveResult(success, fail, warn, del []iface.StoreItem) error {
	if len(success) < 1 && len(fail) < 1 && len(warn) < 1 && len(del) < 1 {
		return nil
	}

	cmd := iproto.CusCmd_ReportLiveResult
	req := &iproto.LiveResultRequest{
		Success: StoreItemToUint64(success...),
		Fail:    StoreItemToUint64(fail...),
		Warn:    StoreItemToUint64(warn...),
		Del:     StoreItemToUint64(del...),
	}

	err := p.CallLeaderApply(cmd, req)
	if err != nil {
		logs.Error("doPublishLiveResult() publish msg result error.", err)
		return err
	}
	//
	p.publishQueue.Delete(success...) // 成功的删除
	p.publishQueue.Delete(fail...)    // 失败的放入死信队列
	p.publishQueue.Delete(warn...)    // 警告数据删除，
	p.publishQueue.Delete(del...)     // 异常数据删除
	//
	p.deadQueue.StoreBatch(fail...)
	//
	return nil
}

func (p *Manager) doPublishDeadResult(success, fail, warn, del []iface.StoreItem) error {
	if len(success) < 1 && len(fail) < 1 && len(warn) < 1 && len(del) < 1 {
		return nil
	}

	cmd := iproto.CusCmd_ReportDeadResult
	req := &iproto.DeadResultRequest{
		Success: StoreItemToUint64(success...),
		Fail:    StoreItemToStorePack(fail...),
		Warn:    StoreItemToUint64(warn...),
		Del:     StoreItemToUint64(del...),
	}
	err := p.CallLeaderApply(cmd, req)
	if err != nil {
		logs.Error("publish msg result error.", err)
		return err
	}
	//
	p.deadQueue.Delete(success...)   // 成功的删除
	p.deadQueue.UpdateBatch(fail...) // 失败的更新retry
	p.deadQueue.Delete(warn...)      // 警告数据删除，
	p.deadQueue.Delete(del...)       // 异常数据删除
	//
	return nil
}

// follower 尝试重新发送死信消息
func (p *Manager) doRepublishMsg(list ...iface.StoreItem) error {
	if len(list) < 1 {
		return nil
	}
	// 发布成功的数据列表
	okList := make([]iface.StoreItem, 0, len(list))
	// 发布失败的数据列表，需要更新retry
	deadList := make([]iface.StoreItem, 0, len(list))
	// 过期未发布的数据列表
	warnList := make([]iface.StoreItem, 0, len(list))
	// 异常数据，要删除列表
	delList := make([]iface.StoreItem, 0, len(list))
	// node := p.CurrentServerMeta.RaftAddr() // 当前节点

	for _, it := range list {
		// var infos []iface.IConsumer
		var booked bool
		if id := it.GetCmdId(); id > 0 {
			_, booked = p.subscribeMap.GetConsumer(id)
		} else {
			logs.Warn("doRepublishMsg() cmdid empty")
			continue
		}

		if !booked {
			// 过期未订阅，未发布的数据
			if /*it.MaxRetry() < 1 &&*/ it.HasExpire(p.cfg.DeadDelay.Seconds() * 2) {
				warnList = append(warnList, it)
			}
			logs.Warn("doPublishMsg() cmdid unbooked=[%v]", it)
			// 再等2轮，等其他节点发布成功，然后删除
			continue
		}

		infos := p.subscribeMap.GetConsumerByKey(it.GetPeers()...)

		err := p.publish(infos, it, true)
		if err != nil {
			deadList = append(deadList, it)
			logs.Error("doPublishMsg() got err=[%v]", err)
		} else {
			okList = append(okList, it)
		}
	}

	p.doPublishDeadResult(okList, deadList, warnList, delList)
	return nil
}

// 上报订阅偏差
func (p *Manager) doReportSubscribeInfoOffset(act bool, cmdIdList []int, key string) {
	req := &iproto.SyncSubscribeInfoOffsetRequest{
		Action: act,
		Ids:    util.IntToUint32(cmdIdList...),
		Key:    key,
	}
	logs.Debug("doReportSubscribeInfoOffset() args=[%v]", req)

	cmd := iproto.CusCmd_ReportSubscribeInfoOffset
	err := p.CallLeaderApply(cmd, req)
	if err != nil {
		logs.Error("doReportSubscribeInfoOffset() got err. %v", err)
	}
}

// 发送给连接
func (p *Manager) publish(infos []iface.IConsumer, msg iface.StoreItem, dead bool) (err error) {
	for _, info := range infos {
		if e := info.Send(msg); e != nil {
			err = e
			msg.SetRecordTime(time.Now().Unix())
			msg.PutPeers(info.GetKey()) // 保存失败的消费者key, 同时retry自增
		} else {
			msg.DelPeers(info.GetKey()) // 发送成功了，则删除peers
		}
	}
	return
}

// leader同步集群里所有的订阅信息
func (p *Manager) syncSubscribedInfo(data map[uint32]uint32) {
	if !p.imLeader {
		return
	}
	req := &iproto.SyncSubscribeInfoRequest{
		Data: data,
	}
	err := p.CallLeaderApply(iproto.CusCmd_SyncSubscribeInfo, req)
	if err != nil {
		logs.Error("syncSubscribedInfo() got err.", err)
	}
}

// 启动的时候，同步一次已订阅的信息列表
func (p *Manager) doAskSubscribedInfo() {
	if !p.imLeader && !p.imFollower {
		logs.Debug("doAskSubscribedInfo() cadidator action")
		return
	}
	if atomic.CompareAndSwapInt32(&p.syncClusterSubscribedInfoOnStart, 0, 1) {
		cmd := iproto.CusCmd_AskSubscribedInfo.String()
		req := &iproto.AskSubscribedInfoRequest{}
		factory := func() proto.Message {
			return &iproto.AskSubscribedInfoReply{}
		}
		// leader广播搜集
		if p.imLeader {
			logs.Debug("doAskSubscribedInfo() leader action")
			out, err := p.ClusterRpcCall(cmd, req, factory, false, p.cfg.Retry, true)
			if err != nil {
				logs.Error("doAskSubscribedInfo() got err. %v", err)
			}
			//
			tmp := map[uint32]uint32{}
			for _, v := range out {
				it, ok := v.(*iproto.AskSubscribedInfoReply)
				if !ok {
					continue
				}
				for k, v := range it.Data {
					if mv, ok := tmp[k]; ok {
						if v > mv {
							tmp[k] = v
						}
					} else {
						tmp[k] = v
					}
				}
			}
			if len(tmp) > 0 {
				p.allSubscribeMap.Set(tmp)
				//
				p.syncSubscribedInfo(tmp)
			}
		} else {
			logs.Debug("doAskSubscribedInfo() follower action")
			// follower向leader请求
			cmd := iproto.CusCmd_AskSubscribedInfo.String()
			addr := p.CurrentLeaderMeta.RpcAddr()

			got, err := p.RpcCall(addr, cmd, p.IncreasedAskId(), req)
			if err != nil {
				logs.Error("doAskSubscribedInfo() got err=%v", err)
				return
			}
			rep := &iproto.AskSubscribedInfoReply{}
			if err = ptypes.UnmarshalAny(got.Data, rep); err != nil {
				logs.Error("doAskSubscribedInfo() got err=%v", err)
				return
			}
			//
			p.allSubscribeMap.Set(rep.Data)
		}

	}
}
