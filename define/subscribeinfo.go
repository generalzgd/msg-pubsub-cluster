/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: subscribeinfo.go
 * @time: 2020/1/2 11:28 上午
 * @project: msgsubscribesvr
 */

package define

import (
	`fmt`

	`github.com/generalzgd/msg-subscriber/iface`
)

const (
	SubscribeTypeTcp iface.SubscribeType = 0
	SubscribeTypeRpc iface.SubscribeType = 1
)

// 订阅的消费者信息
type SubscribeInfo struct {
	FromType    iface.SubscribeType
	ConsumerKey string // 消费者唯一key, 消费者订阅的时候，需要提供唯一的key
	Sender      iface.ISender
}

func (p *SubscribeInfo) GetConnId() uint32 {
	if p.Sender == nil {
		return 0
	}
	return p.Sender.GetConnId()
}

func (p *SubscribeInfo) GetType() iface.SubscribeType {
	return p.FromType
}

func (p *SubscribeInfo) GetAddress() string {
	if p.Sender == nil {
		return ""
	}
	return p.Sender.GetAddress()
}

func (p *SubscribeInfo) Equal(tar iface.IConsumer) bool {
	if p.FromType != tar.GetType() {
		return false
	}
	if p.ConsumerKey == tar.GetKey() {
		return true
	}
	return false
}

func (p *SubscribeInfo) GetKey() string {
	if len(p.ConsumerKey) > 0 {
		return p.ConsumerKey
	}
	return fmt.Sprintf("type:%d>%s", p.FromType, p.GetAddress())
}

func (p *SubscribeInfo) Send(item iface.StoreItem) error {
	if p.Sender != nil {
		return p.Sender.SendItem(item)
	}
	return SenderEmpty
}
