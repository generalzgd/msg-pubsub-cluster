/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: post_handler_test.go.go
 * @time: 2020/1/3 4:16 下午
 * @project: msgsubscribesvr
 */

package mgr

import (
	`encoding/json`
	`os`
	`testing`
	`time`

	`github.com/astaxie/beego/logs`
	`github.com/generalzgd/svr-config/ymlcfg`

	`github.com/generalzgd/msg-subscriber/codec`
	`github.com/generalzgd/msg-subscriber/config`
	`github.com/generalzgd/msg-subscriber/define`
	`github.com/generalzgd/msg-subscriber/iface`
	`github.com/generalzgd/msg-subscriber/storage`
)

var (
	m   = NewManager()
	cfg = config.AppConfig{
		Name:          "",
		LogLevel:      7,
		TreeDegree:    8,
		MaxPackDelay:  time.Second,
		CleanDelay:    time.Second * 60,
		Retry:         1,
		QueueMode:     storage.ModeInMem,
		BatchSize:     50,
		DeadBatchSize: 100,
		DeadDelay:     time.Second * 30,
		Consul:        ymlcfg.ConsulConfig{},
		PostCfg: ymlcfg.TcpCfg{
			Port: 17048,
		},
		SubscribeCfg: config.SubscribeCfg{
			Port: 17049,
		},
		Cluster: ymlcfg.ClusterConfig{
			NodeType: 0,
			SerfPort: 7000,
			RaftPort: 7001,
			RpcPort:  7002,
			HttpPort: 7003,
			Except:   1,
			Peers:    []string{},
		},
		GrpcPool: ymlcfg.ConnPool{
			InitNum:     1,
			CapNum: 5,
		},
		Decode: config.DecodeCfg{
			HeadSize: 24,
			CmdPos:   22,
			CmdSize:  2,
		},
	}
)

func init() {
	err := m.Init(cfg)
	if err != nil {
		logs.Error("Init() got err=(%v)", err)
		os.Exit(1)
	}
	m.StartTcpWork()
}

func TestManager_onSubscribeHandler(t *testing.T) {
	req := define.SubscribeReq{
		ConsumerKey: "a",
		CmdId:       []uint16{1, 2, 3},
	}
	bts, _ := json.Marshal(&req)
	postPack := define.PostPacket{
		Length: uint32(len(bts)),
		CmdId:  7681,
		Body:   bts,
	}
	pk := codec.NewDataPack(postPack.Serialize(), cmdDecoder, bodyDecoder)

	type args struct {
		conn iface.IPackSender
		pack *codec.DataPack
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "TestManager_onSubscribeHandler",
			args: args{
				conn: nil,
				pack: pk,
			},
		},
	}
	time.Sleep(time.Second * 3)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.onSubscribeHandler(tt.args.conn, tt.args.pack)
		})
	}
}

func TestManager_onUnsubscribeHandler(t *testing.T) {
	req := define.SubscribeReq{
		ConsumerKey: "a",
		CmdId:       []uint16{1, 2, 3},
	}
	bts, _ := json.Marshal(&req)
	postPack := define.PostPacket{
		Length: uint32(len(bts)),
		CmdId:  7681,
		Body:   bts,
	}
	pk := codec.NewDataPack(postPack.Serialize(), cmdDecoder, bodyDecoder)

	type args struct {
		conn iface.IPackSender
		pack *codec.DataPack
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "TestManager_onUnsubscribeHandler",
			args: args{
				conn: nil,
				pack: pk,
			},
		},
	}

	m.subscribeMap.Put([]int{2, 3, 4}, &define.SubscribeInfo{
		ConsumerKey: "a",
	})
	time.Sleep(time.Second * 3)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.onUnsubscribeHandler(tt.args.conn, tt.args.pack)
		})
	}
}

func TestManager_onPostMsgHandler(t *testing.T) {

	req := struct {
		CmdId   string `json:"cmdid"`
		Content string `json:"content"`
	}{
		CmdId:   "chatmessage",
		Content: "hahahah",
	}
	bts, _ := json.Marshal(&req)
	postPack := define.PostPacket{
		Length: uint32(len(bts)),
		CmdId:  7681,
		Body:   bts,
	}
	pk := codec.NewDataPack(postPack.Serialize(), cmdDecoder, bodyDecoder)

	type args struct {
		conn iface.IPackSender
		pack *codec.DataPack
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "TestManager_onPostMsgHandler",
			args: args{
				conn: nil,
				pack: pk,
			},
		},
	}
	time.Sleep(time.Second * 3)
	m.doneSubscribe([]int{3846}, &define.SubscribeInfo{
		ConsumerKey: "a",
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.onTcpMsgHandler(tt.args.conn, tt.args.pack)
		})
	}
}
