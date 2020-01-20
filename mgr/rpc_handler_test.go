/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: rpc_handler_test.go.go
 * @time: 2020/1/17 3:01 下午
 * @project: msgsubscribesvr
 */

package mgr

import (
	`context`
	`testing`
	`time`

	`github.com/generalzgd/cluster-plugin/plugin`
	`github.com/golang/protobuf/ptypes`

	`github.com/generalzgd/msg-subscriber/iproto`
)

func TestManager_onRaftApply(t *testing.T) {
	req := &iproto.SyncSubscribeInfoOffsetRequest{
		Action: false,
		Ids:    []uint32{234, 1235},
		Key:    "asdf",
	}
	an, _ := ptypes.MarshalAny(req)
	l := &iproto.ApplyLogEntry{
		Cmd:  iproto.CusCmd_ReportSubscribeInfoOffset,
		Data: an,
	}
	anl, _ := ptypes.MarshalAny(l)

	callReq := &plugin.CallRequest{
		Cmd:  iproto.CusCmd_RaftApply.String(),
		Id:   1,
		Data: anl,
	}

	type args struct {
		ctx context.Context
		req *plugin.CallRequest
	}
	tests := []struct {
		name    string
		args    args
		wantRep *plugin.CallReply
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "TestManager_onRaftApply",
			args: args{
				ctx: context.Background(),
				req: callReq,
			},
		},
	}
	time.Sleep(time.Second * 3)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gotRep, err := m.onRaftApply(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("onRaftApply() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("onRaftApply() gotRep = %v", gotRep)
		})
	}
}
