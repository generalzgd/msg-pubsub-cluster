/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: flowpack_test.go.go
 * @time: 2020/1/3 11:44 上午
 * @project: msgsubscribesvr
 */

package define

import (
	`encoding/binary`
	`encoding/json`
	`sync`
	`testing`

	`github.com/generalzgd/msg-subscriber/codec`
	`github.com/generalzgd/msg-subscriber/codec/body`
	`github.com/generalzgd/msg-subscriber/codec/head`
	`github.com/generalzgd/msg-subscriber/iproto`
)

func TestFlowPack_ToStorePack(t *testing.T) {
	obj := struct {
		Cmdid string `json:"cmdid"`
	}{
		Cmdid: "publish",
	}
	by, _ := json.Marshal(obj) // "{\"cmdid\":\"publish\"}"
	postPk := PostPacket{
		Length: uint32(len(by)),
		CmdId:  7685,
		ToType: 100,
		ToIp:   0,
		Body:   by,
	}
	cmdDecoder := head.NewDecoder(22, 2, 0, 4, binary.LittleEndian)
	bodyDecoder := body.NewDecoder(24, binary.LittleEndian)
	pk := codec.NewDataPack(postPk.Serialize(), cmdDecoder, bodyDecoder)

	type fields struct {
		Index      uint64
		peersLock  sync.RWMutex
		Peers      map[string]uint32
		Pack       *codec.DataPack
		RecordTime int64
	}
	tests := []struct {
		name   string
		fields fields
		want   iproto.StorePack
	}{
		// TODO: Add test cases.
		{
			name: "TestFlowPack_ToStorePack",
			fields: fields{
				Index: 11,
				Peers: map[string]uint32{
					"booker": 1,
				},
				Pack: pk,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &FlowPack{
				Index:      tt.fields.Index,
				peersLock:  tt.fields.peersLock,
				Peers:      tt.fields.Peers,
				Pack:       tt.fields.Pack,
				RecordTime: tt.fields.RecordTime,
			}
			t.Logf("GetCmdStr() = %v", p.GetCmdStr())
			got := p.ToStorePack()
			t.Logf("ToStorePack() = %v", got)
			//
			to := NewFlowPack(cmdDecoder, bodyDecoder)
			to.FromStorePack(got)
			t.Logf("GetCmdStr() = %v, from = %v", to.GetCmdStr(), to.Pack.String())
		})
	}
}
