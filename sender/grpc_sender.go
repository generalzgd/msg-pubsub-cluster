/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: sender.go
 * @time: 2020/1/17 3:56 下午
 * @project: msgsubscribesvr
 */

package sender

import (
	`context`
	`strings`
	`time`

	`github.com/astaxie/beego/logs`
	ctrl `github.com/generalzgd/grpc-svr-frame/grpc-ctrl`
	`github.com/generalzgd/svr-config/ymlcfg`
	grpcpool `github.com/processout/grpc-go-pool`

	`github.com/generalzgd/msg-subscriber/define`
	`github.com/generalzgd/msg-subscriber/iface`
	`github.com/generalzgd/msg-subscriber/iproto`
)

func NewGrpcSender(svrName, addr, clientAddr string, ctrl *ctrl.GrpcController) *GrpcSender {
	return &GrpcSender{
		svrName:   svrName,
		address:   addr,
		rpcClient: clientAddr,
		grpcCtrl:  ctrl,
	}
}

type GrpcSender struct {
	svrName   string               // 要发布的，rpc服务名
	address   string               // 要发布的，rpc地址/或tcp地址, 含端口
	rpcClient string               // 要发布的，消费者rpc 对象地址
	grpcCtrl  *ctrl.GrpcController //
}

func (p *GrpcSender) GetConnId() uint32 {
	return 0
}

func (p *GrpcSender) GetAddress() string {
	return p.address
}

func (p *GrpcSender) SendItem(item iface.StoreItem) error {
	if msg, ok := item.(*define.FlowPack); ok && p.grpcCtrl != nil {
		ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
		var conn *grpcpool.ClientConn
		var err error

		if strings.HasPrefix(p.rpcClient, "consul:///") {
			conn, err = p.grpcCtrl.GetGrpcConnWithLB(ymlcfg.EndpointConfig{
				Name:    p.svrName,
				Address: p.rpcClient,
			}, ctx)
		} else {
			conn, err = p.grpcCtrl.GetGrpcConn(p.svrName, p.rpcClient, ymlcfg.EndpointConfig{}, ctx)
		}
		if err == nil {
			defer conn.Close()
			//
			client := iproto.NewConsumerClient(conn.ClientConn)
			_, err = client.Publish(ctx, &iproto.PublishRequest{
				Index: msg.Index,
				Data:  msg.Pack.GetData(),
			})
			if err != nil {
				logs.Error("SubscribeInfo Send() got err: %v", err)
			}
		} else {
			logs.Error("SubscribeInfo Send() got err: %v", err)
		}
		return err
	}
	return define.SendDataErr
}
