/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: raft_handler_test.go.go
 * @time: 2020/1/15 2:11 下午
 * @project: msgsubscribesvr
 */

package mgr

import (
	`testing`

	`github.com/golang/protobuf/proto`
	`github.com/golang/protobuf/ptypes`
	`github.com/golang/protobuf/ptypes/any`

	`github.com/generalzgd/msg-subscriber/iproto`
)

func TestManager_onRaftSyncClusterIndex(t *testing.T) {
	req := &iproto.SyncClusterIndexRequest{Index: 1234}
	an, _ := ptypes.MarshalAny(req)

	type args struct {
		data *any.Any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "TestManager_onRaftSyncClusterIndex",
			args: args{data: an},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := m.onRaftSyncClusterIndex(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("onRaftSyncClusterIndex() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_onRaftSubscribeInfoOffset(t *testing.T) {

	req := &iproto.SyncSubscribeInfoOffsetRequest{
		Action: false,
		Ids:    []uint32{1, 2, 3},
		Key:    "aa",
	}
	an, err := ptypes.MarshalAny(req)
	if err != nil {
		return
	}
	l := &iproto.ApplyLogEntry{
		Cmd:                  iproto.CusCmd_ReportSubscribeInfoOffset,
		Data:                 an,
	}
	bts, err := proto.Marshal(l)
	if err != nil {
		return
	}

	l2 := &iproto.ApplyLogEntry{}
	if err := proto.Unmarshal(bts, l2); err != nil {
		return
	}

	type args struct {
		data *any.Any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "TestManager_onRaftSubscribeInfoOffset",
			args: args{data: l2.Data},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := m.onRaftSubscribeInfoOffset(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("onRaftSubscribeInfoOffset() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
