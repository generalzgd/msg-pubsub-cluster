/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: postmsg.go
 * @time: 2019/12/24 4:12 下午
 * @project: packagesubscribesvr
 */

package define

type SubscribeReq struct {
	SvrName string `json:"svrname"`
	// 废弃字符串的cmd
	//CmdStr []string `json:"cmds"`
	ConsumerKey string   `json:"consumerkey"`
	CmdId       []uint16 `json:"ids"`
}

type SubscribeAck struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type UnsubscribeReq struct {
	SvrName string `json:"svrname"`
	// 废弃字符串的cmd
	//CmdStr []string `json:"cmds"`
	ConsumerKey string   `json:"consumerkey"`
	CmdId       []uint16 `json:"ids"`
}

type UnsubscribeAck struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type PublishPack struct {
}
