/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: err.go
 * @time: 2019/12/23 7:38 下午
 * @project: packagesubscribesvr
 */

package define

import (
	`errors`
)

var (
	LeaderNot = errors.New("im not leader")
	LeaderNow = errors.New("im leader")
	BucketEmpty = errors.New("bucket empty")
	StorageEmpty = errors.New("target store empty")
	ErrorCode = errors.New("got error code")
	DataEmpty = errors.New("data empty")
	ParamNil = errors.New("param nil")
	ApplyErr = errors.New("leader apply fail")
	IndexErr = errors.New("get cluster increased index fail")
	PostSendErr = errors.New("post link send fail")
	SendDataErr = errors.New("send data is not *FlowPack")
	SenderEmpty = errors.New("sender empty")
	//CrossBorder = errors.New("cross the border")
	//ProtoDecodeErr = errors.New("proto decode error")
	//ProtoEncodeErr = errors.New("proto encode error")
)
