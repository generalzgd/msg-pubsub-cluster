/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: isender.go
 * @time: 2020/1/17 3:56 下午
 * @project: msgsubscribesvr
 */

package iface

import (
	gotcp `github.com/generalzgd/securegotcp`
)

type IPackSender interface {
	Send(packet gotcp.Packet)error
	GetAddress() string
	GetConnId() uint32
}

type ISender interface {
	SendItem(StoreItem)error
	GetAddress() string
	GetConnId() uint32
}
