/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: rpc_handler.go
 * @time: 2019/12/23 6:46 下午
 * @project: packagesubscribesvr
 */

package mgr

import (
	`context`

	`github.com/astaxie/beego/logs`
	`github.com/generalzgd/cluster-plugin/plugin`
	`github.com/golang/protobuf/ptypes`

	`github.com/generalzgd/msg-subscriber/define`
	`github.com/generalzgd/msg-subscriber/iproto`
)

type RpcHandler func(context.Context, *plugin.CallRequest) (*plugin.CallReply, error)

func (p *Manager) initRpcHandlers() {
	p.rpcHandlers = map[string]RpcHandler{
		iproto.CusCmd_RaftApply.String():           p.onRaftApply,
		iproto.CusCmd_CollectClusterIndex.String(): p.onRpcCollectClusterIndex,
		iproto.CusCmd_AskIndex.String():            p.onRpcAskIndex,
		iproto.CusCmd_AskIncreasedIndex.String():   p.onRpcAskIncreasedIndex,
		iproto.CusCmd_AskSubscribedInfo.String():   p.onRpcAskSubscribedInfo,
	}
}

// leader收到集群follower节点rpc请求raft同步
func (p *Manager) onRaftApply(ctx context.Context, req *plugin.CallRequest) (rep *plugin.CallReply, err error) {
	rep = &plugin.CallReply{
		Cmd: req.Cmd,
		Id:  req.Id,
	}
	if !p.imLeader {
		err = define.LeaderNot
		return
	}

	// RaftApplyRequest
	args := &iproto.ApplyLogEntry{}
	if err = ptypes.UnmarshalAny(req.Data, args); err != nil {
		logs.Debug("onRaftApply() got err=[%v]", err)
		return
	} else {
		// logs.Debug("onRaftApply() got req=[%v]", args.Data)
	}

	err = p.CallLeaderApply(args.Cmd, args.Data)
	return
}

// 响应自己的索引值
func (p *Manager) onRpcCollectClusterIndex(ctx context.Context, req *plugin.CallRequest) (rep *plugin.CallReply, err error) {
	// logs.Debug("onRpcCollectClusterIndex() got=%v", req)
	rep = &plugin.CallReply{
		Cmd: req.Cmd,
		Id:  req.Id,
	}

	obj := &iproto.CollectClusterIndexReply{
		Res:   &iproto.CommonReply{},
		Index: p.GetIndex(), // p.getIncreasedIdx(),
	}
	buf, err := ptypes.MarshalAny(obj)
	if err != nil {
		logs.Error("onRpcCollectClusterIndex() got err=%v", err)
		return
	}
	// logs.Debug("onRpcCollectClusterIndex() return=%v", obj)
	rep.Data = buf
	return
}

// leader 返回最新的索引
func (p *Manager) onRpcAskIndex(ctx context.Context, req *plugin.CallRequest) (rep *plugin.CallReply, err error) {
	rep = &plugin.CallReply{
		Cmd: req.Cmd,
		Id:  req.Id,
	}

	args := &iproto.AskIndexRequest{}
	if err = ptypes.UnmarshalAny(req.Data, args); err != nil {
		logs.Error("onRpcAskIndex() got err=%v", err)
		return
	} else {
		// logs.Debug("onRpcAskIndex() got req=%v", args)
	}
	// 同步索引
	if args.Index > p.GetIndex() {
		p.SetIndex(args.Index)
		p.syncClusterIndex(args.Index)
	}

	out := &iproto.AskIndexReply{
		Res:   &iproto.CommonReply{},
		Index: p.GetIndex(),
	}
	bts, err := ptypes.MarshalAny(out)
	if err != nil {
		logs.Error("onRpcAskIndex() got err=%v", err)
		return
	}
	// logs.Debug("onRpcAskIndex() return=%v", out)
	rep.Data = bts
	return
}

// leader返回自增加的index，并同步给其他节点
func (p *Manager) onRpcAskIncreasedIndex(ctx context.Context, req *plugin.CallRequest) (rep *plugin.CallReply, err error) {
	rep = &plugin.CallReply{
		Cmd: req.Cmd,
		Id:  req.Id,
	}
	index, _ := p.IncreaseIndex()
	cmd := iproto.CusCmd_SyncClusterIndex
	if err = p.CallLeaderApply(cmd, &iproto.SyncClusterIndexRequest{
		Index: index,
	}); err != nil {
		logs.Error("onRpcAskIncreasedIndex() got err=%v", err)
		return
	}

	out := &iproto.AskIncreasedIndexReply{
		Res:   &iproto.CommonReply{},
		Index: index,
	}
	bts, err := ptypes.MarshalAny(out)
	if err != nil {
		logs.Error("onRpcAskIncreasedIndex() got err=%v", err)
		return
	}
	// logs.Debug("onRpcAskIncreasedIndex() return=%v", out)
	rep.Data = bts
	return
}

// leader/follower响应已订阅的信息
func (p *Manager) onRpcAskSubscribedInfo(ctx context.Context, req *plugin.CallRequest) (rep *plugin.CallReply, err error) {
	rep = &plugin.CallReply{
		Cmd: req.Cmd,
		Id:  req.Id,
	}
	args := &iproto.AskSubscribedInfoRequest{}
	if err = ptypes.UnmarshalAny(req.Data, args); err != nil {
		logs.Error("onRpcAskSubscribedInfo() got err=%v", err)
		return
	} else {
		// logs.Debug("onRpcAskSubscribedInfo() got req=%v", args)
	}
	p.allSubscribeMap.Set(args.Data)
	//
	out := &iproto.AskSubscribedInfoReply{
		Res:  &iproto.CommonReply{},
		Data: p.allSubscribeMap.Copy(),
	}
	bts, err := ptypes.MarshalAny(out)
	if err != nil {
		logs.Error("onRpcAskSubscribedInfo() got err=%v", err)
		return
	}
	// logs.Debug("onRpcAskSubscribedInfo() return=%v", out)
	rep.Data = bts
	return
}
