/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: isonsumer.go
 * @time: 2020/1/2 2:11 下午
 * @project: msgsubscribesvr
 */

package iface

type SubscribeType int

// 消费者
type IConsumer interface {
	GetType() SubscribeType
	GetAddress() string
	Equal(consumer IConsumer) bool
	GetKey() string
	GetConnId() uint32
	Send(item StoreItem) error
}
