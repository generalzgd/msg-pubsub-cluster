/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: raft_handler.go
 * @time: 2020/1/10 6:04 下午
 * @project: msgsubscribesvr
 */

package mgr

import (
	`github.com/astaxie/beego/logs`
	`github.com/golang/protobuf/ptypes`
	`github.com/golang/protobuf/ptypes/any`

	`github.com/generalzgd/msg-subscriber/iproto`
	`github.com/generalzgd/msg-subscriber/util`
)

type RaftHandler func(data *any.Any) error

func (p *Manager) initRaftHandlers() {
	p.raftHandlers = map[iproto.CusCmd]RaftHandler{
		iproto.CusCmd_SyncClusterIndex:          p.onRaftSyncClusterIndex,
		iproto.CusCmd_SyncSubscribeInfo:         p.onRaftSyncSubscribeInfo,
		iproto.CusCmd_ReportNewMsg:              p.onRaftReportNewMsg,
		iproto.CusCmd_ReportLiveResult:          p.onRaftReportLiveResult,
		iproto.CusCmd_ReportDeadResult:          p.onRaftReportDeadResult,
		iproto.CusCmd_ReportSubscribeInfoOffset: p.onRaftSubscribeInfoOffset,
	}
}

// raft 同步索引
func (p *Manager) onRaftSyncClusterIndex(data *any.Any) error {
	// logs.Debug("onRaftSyncClusterIndex()")

	args := &iproto.SyncClusterIndexRequest{}
	if err := ptypes.UnmarshalAny(data, args); err != nil {
		logs.Error("onRaftSyncClusterIndex() got err=%v", err)
		return err
	} else {
		logs.Debug("onRaftSyncClusterIndex() got args=%v", args.String())
	}
	p.SetIndex(args.Index)
	return nil
}

// raft 同步集群里订阅信息
func (p *Manager) onRaftSyncSubscribeInfo(data *any.Any) error {
	// logs.Debug("onRaftSyncSubscribeInfo()")

	args := &iproto.SyncSubscribeInfoRequest{}
	if err := ptypes.UnmarshalAny(data, args); err != nil {
		logs.Error("onRaftSyncSubscribeInfo() got err=%v", err)
		return err
	} else {
		logs.Debug("onRaftSyncSubscribeInfo() got args=%v", args.String())
	}
	p.allSubscribeMap.Set(args.Data)
	return nil
}

// 收到同步的新消息
func (p *Manager) onRaftReportNewMsg(data *any.Any) error {
	// logs.Debug("onRaftReportNewMsg() got=[%v]", string(data))

	args := &iproto.NewMsgRequest{}
	if err := ptypes.UnmarshalAny(data, args); err != nil {
		logs.Error("onRaftReportNewMsg() got err=[%v]", err)
		return err
	} else {
		logs.Debug("onRaftReportNewMsg() got args=[%v, %v]", args.Data, len(args.Data))
	}
	if len(args.Data) > 0 {
		item := StorePackToStoreItem(args.Data...)
		if err := p.publishQueue.StoreBatch(item...); err != nil {
			logs.Error("onRaftReportNewMsg() got err=[%v]", err)
			return err
		}
	}
	return nil
}

// raft 同步消息发布结果
func (p *Manager) onRaftReportLiveResult(data *any.Any) error {
	// logs.Debug("onRaftReportLiveResult()")

	args := &iproto.LiveResultRequest{}
	if err := ptypes.UnmarshalAny(data, args); err != nil {
		logs.Error("onRaftReportLiveResult() got err=%v", err)
		return err
	} else {
		logs.Debug("onRaftReportLiveResult() got args=%v", args.String())
	}
	//
	p.publishQueue.DeleteByIndexes(args.Success...)
	fail, _ := p.publishQueue.DeleteByIndexes(args.Fail...)
	p.publishQueue.DeleteByIndexes(args.Warn...)
	p.publishQueue.DeleteByIndexes(args.Del...)
	//
	//fail := StorePackToStoreItem(args.Fail...)
	//p.publishQueue.Delete(StorePackToStoreItem(args.Success...)...) // 成功的删除
	//p.publishQueue.Delete(fail...)                                  // 失败的放入死信队列
	//p.publishQueue.Delete(StorePackToStoreItem(args.Warn...)...)    // 警告数据删除，
	//p.publishQueue.Delete(StorePackToStoreItem(args.Del...)...)     // 异常数据删除
	p.deadQueue.StoreBatch(fail...)
	return nil
}

// raft 同步死信消息发布结果
func (p *Manager) onRaftReportDeadResult(data *any.Any) error {
	// logs.Debug("onRaftReportDeadResult()")

	args := &iproto.DeadResultRequest{}
	if err := ptypes.UnmarshalAny(data, args); err != nil {
		logs.Error("onRaftReportDeadResult() got err=%v", err)
		return err
	} else {
		logs.Debug("onRaftReportDeadResult() got args=%v", args.String())
	}
	//
	p.deadQueue.DeleteByIndexes(args.Success...)
	p.deadQueue.UpdateBatch(StorePackToStoreItem(args.Fail...)...)
	p.deadQueue.DeleteByIndexes(args.Warn...)
	p.deadQueue.DeleteByIndexes(args.Del...)
	//
	//p.deadQueue.Delete(StorePackToStoreItem(args.Success...)...)   // 成功的删除
	//p.deadQueue.UpdateBatch(StorePackToStoreItem(args.Fail...)...) // 失败的更新retry
	//p.deadQueue.Delete(StorePackToStoreItem(args.Warn...)...)      // 警告数据删除，
	//p.deadQueue.Delete(StorePackToStoreItem(args.Del...)...)       // 异常数据删除
	return nil
}

// raft 同步订阅信息
func (p *Manager) onRaftSubscribeInfoOffset(data *any.Any) error {
	// logs.Debug("onRaftSubscribeInfoOffset()")

	args := &iproto.SyncSubscribeInfoOffsetRequest{}
	if err := ptypes.UnmarshalAny(data, args); err != nil {
		logs.Error("onRaftSubscribeInfoOffset() got err=[%v]", err)
		return err
	} else {
		logs.Debug("onRaftSubscribeInfoOffset() got args=[%v, %v, %v]", args.Action, args.Ids, args.Key)
	}
	if len(args.Ids) > 0 {
		if args.Action {
			p.allSubscribeMap.Put(util.Uint32ToInt(args.Ids...)...)
		} else {
			p.allSubscribeMap.Del(util.Uint32ToInt(args.Ids...)...)
		}
	}
	return nil
}
