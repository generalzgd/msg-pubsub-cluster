/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: subscribeinfo_test.go.go
 * @time: 2020/1/3 2:02 下午
 * @project: msgsubscribesvr
 */

package define

import (
	`encoding/binary`
	`testing`

	`github.com/generalzgd/msg-subscriber/codec/body`
	`github.com/generalzgd/msg-subscriber/codec/head`
	`github.com/generalzgd/msg-subscriber/iface`
)

func TestSubscribeInfo_Send(t *testing.T) {
	cmdDecoder := head.NewDecoder(22, 2, 0, 4, binary.LittleEndian)
	bodyDecoder := body.NewDecoder(24, binary.LittleEndian)
	type fields struct {
		FromType    iface.SubscribeType
		ConsumerKey string
		Sender      iface.ISender
	}
	type args struct {
		item iface.StoreItem
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "TestSubscribeInfo_Send_1",
			fields: fields{
				FromType:    SubscribeTypeTcp,
				ConsumerKey: "abc",
				Sender:      nil, // sender.NewTcpSender("abc", 1, "127.0.0.1:10000", nil),
			},
			args:    args{item: NewFlowPack(cmdDecoder, bodyDecoder)},
			wantErr: true,
		},
		{
			name: "TestSubscribeInfo_Send_2",
			fields: fields{
				FromType:    SubscribeTypeRpc,
				ConsumerKey: "abc",
				Sender:      nil, // sender.NewGrpcSender("abc", "127.0.0.1:10000", "127.0.0.1:10000", nil),
			},
			args:    args{item: NewFlowPack(cmdDecoder, bodyDecoder)},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &SubscribeInfo{
				FromType:    tt.fields.FromType,
				ConsumerKey: tt.fields.ConsumerKey,
				Sender:      tt.fields.Sender,
			}
			if err := p.Send(tt.args.item); (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
