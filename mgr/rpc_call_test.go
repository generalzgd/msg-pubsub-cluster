/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: rpc_call_test.go.go
 * @time: 2020/1/15 12:26 下午
 * @project: msgsubscribesvr
 */

package mgr

import (
	`testing`
	`time`

	`github.com/generalzgd/msg-subscriber/codec`
	`github.com/generalzgd/msg-subscriber/define`
	`github.com/generalzgd/msg-subscriber/iface`
	`github.com/generalzgd/msg-subscriber/iproto`
	`github.com/generalzgd/msg-subscriber/sender`
)

func TestManager_collectClusterIndex(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "TestManager_collectClusterIndex",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut, err := m.collectClusterIndex()
			if (err != nil) != tt.wantErr {
				t.Errorf("collectClusterIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("collectClusterIndex() got = %v", gotOut)
		})
	}
}

func TestManager_syncClusterIndex(t *testing.T) {

	type args struct {
		index uint64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "TestManager_syncClusterIndex",
			args: args{index: 1234},
		},
	}
	time.Sleep(time.Second * 3)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := m.syncClusterIndex(tt.args.index)
			if (err != nil) != tt.wantErr {
				t.Errorf("syncClusterIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestManager_askClusterIndex(t *testing.T) {

	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "TestManager_askClusterIndex",
		},
	}
	time.Sleep(time.Second * 3)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gotOut, err := m.askClusterIndex()
			if (err != nil) != tt.wantErr {
				t.Errorf("askClusterIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("askClusterIndex() gotOut = %v", gotOut)
		})
	}
}

func TestManager_askIncreasedIndex(t *testing.T) {

	tests := []struct {
		name    string
		want    uint64
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "TestManager_askIncreasedIndex",
			want: 1235,
		},
	}
	time.Sleep(time.Second * 3)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := m.askIncreasedIndex()
			if (err != nil) != tt.wantErr {
				t.Errorf("askIncreasedIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("askIncreasedIndex() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_doReportMsg(t *testing.T) {
	contents := []byte(`{"cmdid":"chatmessage", "content":"asdfsadfasdf"}`)
	postPk := define.PostPacket{
		Length: uint32(len(contents)),
		CmdId:  3846,
		Body:   contents,
	}

	pk := codec.NewDataPack(postPk.Serialize(), headDecoder, bodyDecoder)

	type args struct {
		list []iface.StoreItem
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "TestManager_doReportMsg",
			args: args{list: []iface.StoreItem{
				&define.FlowPack{
					Index: 1,
					Pack:  pk,
				},
			}},
		},
	}
	time.Sleep(time.Second * 3)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := m.doReportMsg(tt.args.list...); (err != nil) != tt.wantErr {
				t.Errorf("doReportMsg() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_doSyncClusterIndex(t *testing.T) {

	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{
			name: "TestManager_doSyncClusterIndex",
		},
	}
	time.Sleep(time.Second * 2)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.doSyncClusterIndex()
		})
	}
}

func TestManager_doAskClusterIndex(t *testing.T) {

	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{
			name: "TestManager_doAskClusterIndex",
		},
	}
	time.Sleep(time.Second * 2)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.doAskClusterIndex()
		})
	}
}

func TestManager_doPublishMsg(t *testing.T) {
	contents := []byte(`{"cmdid":"chatmessage", "content":"asdfsadfasdf"}`)
	postPk := define.PostPacket{
		Length: uint32(len(contents)),
		CmdId:  3846,
		Body:   contents,
	}

	pk := codec.NewDataPack(postPk.Serialize(), headDecoder, bodyDecoder)

	type args struct {
		list []iface.StoreItem
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "TestManager_doPublishMsg",
			args: args{list: []iface.StoreItem{
				&define.FlowPack{
					Index: 1,
					Pack:  pk,
				},
			}},
		},
	}
	time.Sleep(time.Second * 2)
	m.doneSubscribe([]int{3846}, &define.SubscribeInfo{
		FromType:    0,
		ConsumerKey: "a",
		Sender:      sender.NewTcpSender("", 1, "127.0.0.1:8000", nil),
	})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := m.doPublishMsg(tt.args.list...); (err != nil) != tt.wantErr {
				t.Errorf("doPublishMsg() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_doRepublishMsg(t *testing.T) {
	contents := []byte(`{"cmdid":"chatmessage", "content":"asdfsadfasdf"}`)
	postPk := define.PostPacket{
		Length: uint32(len(contents)),
		CmdId:  3846,
		Body:   contents,
	}

	pk := codec.NewDataPack(postPk.Serialize(), headDecoder, bodyDecoder)

	type args struct {
		list []iface.StoreItem
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "TestManager_doRepublishMsg",
			args: args{list: []iface.StoreItem{
				&define.FlowPack{
					Index: 1,
					Pack:  pk,
				}},
			},
		},
	}

	time.Sleep(time.Second * 2)
	m.doneSubscribe([]int{3846}, &define.SubscribeInfo{
		FromType:    0,
		ConsumerKey: "a",
		Sender:      sender.NewTcpSender("", 1, "127.0.0.1:8000", nil),
	})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := m.doRepublishMsg(tt.args.list...); (err != nil) != tt.wantErr {
				t.Errorf("doRepublishMsg() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_doReportSubscribeInfoOffset(t *testing.T) {

	type args struct {
		act       bool
		cmdIdList []int
		key       string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "TestManager_doReportSubscribeInfoOffset_1",
			args: args{
				act:       true,
				cmdIdList: []int{3846},
				key:       "a",
			},
		},
		{
			name: "TestManager_doReportSubscribeInfoOffset_2",
			args: args{
				act:       false,
				cmdIdList: []int{3846},
				key:       "a",
			},
		},
	}
	time.Sleep(time.Second * 2)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.doReportSubscribeInfoOffset(tt.args.act, tt.args.cmdIdList, tt.args.key)
		})
	}
}

func TestManager_syncSubscribedInfo(t *testing.T) {

	type args struct {
		data map[string]*iproto.NodeMap
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "TestManager_syncSubscribedInfo",
			args: args{data: map[string]*iproto.NodeMap{
				"node1": {
					Data: map[string]*iproto.IdList{
						"a": {
							Ids: []uint32{3846},
						},
					},
				},
			},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.syncSubscribedInfo(tt.args.data)
		})
	}
}

func TestManager_doAskSubscribedInfo(t *testing.T) {

	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{
			name: "TestManager_doAskSubscribedInfo",
		},
	}
	time.Sleep(time.Second * 2)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.doAskSubscribedInfo()
		})
	}
}
