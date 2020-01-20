/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: manager_test.go.go
 * @time: 2020/1/6 4:05 下午
 * @project: msgsubscribesvr
 */

package mgr

import (
	`encoding/json`
	`testing`
	`time`

	`github.com/generalzgd/msg-subscriber/codec`
	`github.com/generalzgd/msg-subscriber/define`
	`github.com/generalzgd/msg-subscriber/iface`
	`github.com/generalzgd/msg-subscriber/sender`
)

func prepareManagerData(queue iface.IStoreBridge) {
	req := struct {
		CmdId   string `json:"cmdid"`
		Content string `json:"content"`
	}{
		CmdId:   "chatmessage",
		Content: "hahahah1",
	}
	bts, _ := json.Marshal(&req)

	postPk := define.PostPacket{
		Length: uint32(len(bts)),
		CmdId:  3846,
		Body:   bts,
	}

	pk := codec.NewDataPack(postPk.Serialize(), headDecoder, bodyDecoder)

	li := []iface.StoreItem{&define.FlowPack{
		Index: 1,
		Peers: map[string]uint32{
			"a": 1,
		},
		Pack: pk,
	}, &define.FlowPack{
		Index: 20,
		Peers: map[string]uint32{
			"a": 1,
		},
		Pack: pk,
	}, &define.FlowPack{
		Index: 5,
		Peers: map[string]uint32{
			"a": 1,
		},
		Pack: pk,
	}, &define.FlowPack{
		Index: 155,
		Peers: map[string]uint32{
			"a": 1,
		},
		Pack: pk,
	}}

	queue.StoreBatch(li...)

	select {
	case <-m.quit:
		m.quit = make(chan struct{})
	default:
	}
}

func autoQuitGorutineTest() {
	select {
	case <-time.After(time.Minute):
		close(m.quit)
	}
}

/*func TestManager_watchLeaderFlow(t *testing.T) {
	prepareManagerData(m.leaderQueue)

	go m.watchLeaderFlow()

	autoQuitGorutineTest()
}*/

func TestManager_watchPublishFlow(t *testing.T) {
	prepareManagerData(m.publishQueue)
	m.doneSubscribe([]int{3846}, &define.SubscribeInfo{
		FromType:    0,
		ConsumerKey: "a",
		Sender:      sender.NewTcpSender("a", 1, "127.0.0.1", nil),
	})

	go m.watchPublishFlow()

	autoQuitGorutineTest()
}

func TestManager_watchDeadFlow(t *testing.T) {
	prepareManagerData(m.deadQueue)
	m.doneSubscribe([]int{3846}, &define.SubscribeInfo{
		FromType:    0,
		ConsumerKey: "a",
		Sender:      sender.NewTcpSender("a", 1, "127.0.0.1", nil),
	})

	go m.watchDeadFlow()

	autoQuitGorutineTest()
}
