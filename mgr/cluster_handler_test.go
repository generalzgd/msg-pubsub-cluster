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
	"testing"
	`time`

	"github.com/golang/protobuf/proto"

	`github.com/generalzgd/msg-subscriber/define`
	"github.com/generalzgd/msg-subscriber/iproto"
)



func TestManager_CallLeaderApply(t *testing.T) {
	body := `{"inweeklist":1,"speakinroom":0,"toid":0,"extra":"","fromuid":163,"fromid":1652831208,"permission":1,"level":20,"vlevel":4,"slevel":23,"plevel":9,"pos":5,"swline":7,"roomurl":10086,"roomid":874,"gameid":35,"time":1579177976,"swip":3078740009,"chatid":"1579177976204","ava":"1495776830","vdesc":"查看线路","ip":"115.236.48.234","fromname":"你好我好大家好","content":"777","cmdid":"chatmessage"}`
	pk := define.PostPacket{
		Length: uint32(len(body)),
		CmdId:  3846,
		Body:   []byte(body),
	}

	type args struct {
		cusCmd iproto.CusCmd
		msg    proto.Message
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "TestManager_CallLeaderApply",
			args: args{
				cusCmd: iproto.CusCmd_ReportNewMsg,
				msg: &iproto.NewMsgRequest{
					Data: []*iproto.StorePack{
						{
							Index: 1,
							Data:  pk.Serialize(),
						},
						{
							Index: 2,
							Data:  pk.Serialize(),
						},
						{
							Index: 3,
							Data:  pk.Serialize(),
						},
					},
				},
			},
		},
	}
	time.Sleep(time.Second * 2)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := m.CallLeaderApply(tt.args.cusCmd, tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("Manager.CallLeaderApply() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	time.Sleep(time.Second * 15)
}
