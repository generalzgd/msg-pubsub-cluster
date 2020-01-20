/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: postsender.go
 * @time: 2020/1/19 4:37 下午
 * @project: msgsubscribesvr
 */

package sender

import (
	`github.com/generalzgd/msg-subscriber/define`
	`github.com/generalzgd/msg-subscriber/iface`
)

func NewTcpSender(svrName string, connId uint32, addr string, conn iface.IPackSender) *PostSender {
	return &PostSender{
		svrName:   svrName,
		connId:    connId,
		address:   addr,
		tcpClient: conn,
	}
}

type PostSender struct {
	svrName   string           // 要发布的，服务名
	connId    uint32           // 要发布的，tcp socket id，从post发布消息
	address   string           // 要发布的，rpc地址/或tcp地址, 含端口
	tcpClient iface.IPackSender //
}

func (p *PostSender) GetConnId() uint32 {
	return p.connId
}

func (p *PostSender) GetAddress() string {
	return p.address
}

func (p *PostSender) SendItem(item iface.StoreItem) error {
	if p.tcpClient != nil {
		if msg, ok := item.(*define.FlowPack); ok {
			return p.tcpClient.Send(msg.Pack)
		} else {
			return define.SendDataErr
		}
	}
	return define.PostSendErr
}
