/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: clusterhandler.go
 * @time: 2019/12/23 2:54 下午
 * @project: packagesubscribesvr
 */

package mgr

import (
	`context`
	`io`
	`time`

	`github.com/astaxie/beego/logs`
	`github.com/generalzgd/cluster-plugin/plugin`
	`github.com/golang/protobuf/proto`
	`github.com/golang/protobuf/ptypes`
	`github.com/golang/protobuf/ptypes/any`
	`github.com/hashicorp/raft`
	`github.com/hashicorp/serf/serf`

	`github.com/generalzgd/msg-subscriber/define`
	`github.com/generalzgd/msg-subscriber/iproto`
)

func (p *Manager) Persist(sink raft.SnapshotSink) error {
	// panic("implement me")
	/*m := p.allSubscribeMap.Copy()
	bts, err := json.Marshal(m)
	if err != nil {
		sink.Cancel()
		return err
	}
	if _, err := sink.Write(bts); err != nil {
		sink.Cancel()
		return err
	}
	if err := sink.Close(); err != nil {
		sink.Close()
		return err
	}*/
	return nil
}

func (p *Manager) Release() {
}

// *************************************************

// Apply用户数据同步，Leader发起广播给其他节点。不支持收集类操作
func (p *Manager) Apply(l *raft.Log) interface{} {
	// logs.Debug("Apply() got=%v, %s", l, string(l.Data))
	if len(l.Data) > 0 {
		req := &iproto.ApplyLogEntry{}
		if err := proto.Unmarshal(l.Data, req); err == nil {
			if f, ok := p.raftHandlers[req.Cmd]; ok {
				if err := f(req.Data); err == nil {
					return "ok"
				}
			}
		} else {
			logs.Error("Apply() got err=(%v)", err)
		}
	} else {
		logs.Error("Apply() got err=(empty log data)")
	}
	return "fail"
}

func (p *Manager) Snapshot() (raft.FSMSnapshot, error) {
	return p, nil
}

func (p *Manager) Restore(r io.ReadCloser) error {
	logs.Debug("Restore()")

	/*m := map[uint32]uint32{}
	if err := json.NewDecoder(r).Decode(&m); err != nil {
		return err
	}
	p.allSubscribeMap.Set(m)*/
	return nil
}

// **************************************************

// raft
func (p *Manager) OnLeaderSwift(id string) {
	logs.Debug("OnLeaderSwift() got=[%v, %v]", id, p.MemberList())
	if len(id) < 1 {
		p.CurrentLeaderMeta = plugin.ServerMeta{}
		return
	}
	for _, m := range p.MemberList() {
		// logs.Debug("OnLeaderSwift() Addr=[%v]", m.Addr)
		meta := plugin.ServerMeta{
			Name: m.Name,
			Addr: m.Addr,
		}
		if err := meta.FromMap(m.Tags); err != nil {
			continue
		}
		// logs.Debug("OnLeaderSwift() iterator=[%v, %v]", meta, meta.RaftAddr())
		if meta.RaftAddr() == id {
			p.CurrentLeaderMeta = meta
			// logs.Debug("OnLeaderSwift() LeaderMeta=[%v, %v]", meta, meta.RpcAddr())
			break
		}
	}
}

func (p *Manager) OnImLeader() {
	p.imLeader = true
	p.imFollower = false
	// 同步集群索引
	p.doSyncClusterIndex()
	// todo 启动后，同步订阅信息
	p.doAskSubscribedInfo()

	// 尝试启动侦听，接收客户端
	p.tryStartAccept()
	// p.testFsm()
}

/*func (p *Manager) testFsm() {
	r := p.ClusterPlugin.Raft()
	//r.LastIndex()
	//leader := r.Leader()
	f := r.Apply([]byte("hello fsm"), time.Second*5)
	if e := f.Error(); e != nil {
		logs.Error(e)
	}
	logs.Debug("testFsm() got=%v", f.Response())
}*/

func (p *Manager) OnImFollower() {
	p.imLeader = false
	p.imFollower = true
	//
	p.doAskClusterIndex()
	// todo 启动后，同步订阅信息
	p.doAskSubscribedInfo()
	// 尝试启动侦听，接收客户端
	p.tryStartAccept()
}

func (p *Manager) OnNodeReady() {

}

func (p *Manager) OnImCandidate() {
	p.imLeader = false
	p.imFollower = false
}

func (p *Manager) OnImVoter() {
}

func (p *Manager) OnImNonvoter() {
}

func (p *Manager) OnNodeJoin(string) {
}

func (p *Manager) OnNodeLeave(string) {
}

func (p *Manager) OnWarn(serf.UserEvent) {
}

func (p *Manager) OnCustomMsg([]byte) {
}

func (p *Manager) OnQuery(*serf.Query) {
}

func (p *Manager) OnRpcCall(ctx context.Context, req *plugin.CallRequest) (*plugin.CallReply, error) {
	f, ok := p.rpcHandlers[req.Cmd]
	if ok {
		return f(ctx, req)
	}

	return &plugin.CallReply{
		Cmd: req.Cmd,
		Id:  req.Id,
	}, nil
}

// 访问leader执行fsm apply
// @param cuscmd
// @param msg ApplyLogEntry
func (p *Manager) CallLeaderApply(cusCmd iproto.CusCmd, data proto.Message) error {
	l := &iproto.ApplyLogEntry{
		Cmd: cusCmd,
	}

	if an, ok := data.(*any.Any); ok {
		l.Data = an
	} else {
		if an, err := ptypes.MarshalAny(data); err == nil {
			l.Data = an
		} else {
			logs.Error("CallLeaderApply() got err=(%v)", err)
			return err
		}
	}

	// logs.Debug("CallLeaderApply() args=[%v, %v, %v]", p.imLeader, cusCmd, msg)
	if p.imLeader {
		if bts, err := proto.Marshal(l); err == nil {
			f := p.Raft().Apply(bts, time.Second*3)
			if e := f.Error(); e != nil {
				logs.Error("CallLeaderApply err=%v", e)
			}
			r := f.Response()
			if out, ok := r.(string); ok {
				switch out {
				case "fail":
					logs.Error("leader apply fail. %s", cusCmd)
					return define.ApplyErr
				case "ok":
					// logs.Debug("leader apply ok. %s", cusCmd)
					return nil
				}
			}
		}
	} else if p.imFollower {
		addr := p.CurrentLeaderMeta.RpcAddr()
		if _, err := p.ClusterPlugin.RpcCall(addr, iproto.CusCmd_RaftApply.String(), p.IncreasedAskId(), l); err != nil {
			logs.Error("CallLeaderApply() got err=[%v, %v]", err, addr)
			return err
		}
	}
	return nil
}
