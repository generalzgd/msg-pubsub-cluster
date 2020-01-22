/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: datapack.go
 * @time: 2020/1/19 2:53 下午
 * @project: msgsubscribesvr
 */

package codec

import (
	`fmt`

	`github.com/astaxie/beego/logs`

	`github.com/generalzgd/msg-subscriber/codec/cmdstr`
)

//
type DataPack struct {
	cmdId       int
	data        []byte
	headDecoder IHeadCodec
	bodyDecoder IBodyCodec
}

func NewDataPack(data []byte, headDecoder IHeadCodec, bodyDecoder IBodyCodec) *DataPack {
	return &DataPack{
		data:        data,
		headDecoder: headDecoder,
		bodyDecoder: bodyDecoder,
	}
}

func (p *DataPack) Serialize() []byte {
	return p.data
}

func (p *DataPack) GetData() []byte {
	return p.data
}

func (p *DataPack) GetCmdId() int {
	if p.cmdId > 0 {
		return p.cmdId
	}
	val, _, err := p.headDecoder.Read(p.data)
	if err != nil {
		logs.Error("GetCmdId() got err=(%v)", err)
		return 0
	}
	p.cmdId = val
	return val
}

func (p *DataPack) SetHead(cmd, l int) {
	if cmd > 0 {
		err := p.headDecoder.Write(p.data, cmd, l)
		if err != nil {
			logs.Error("SetCmdId() got err=(%v)", err)
		}
	}
}

func (p *DataPack) GetCmdStr() string {
	decoder := cmdstr.NewDecoder()
	val, err := decoder.Read(p.GetPackBody())
	if err != nil {
		return ""
	}
	return val
}

func (p *DataPack) GetPackBody() []byte {
	val, err := p.bodyDecoder.Read(p.data)
	if err != nil {
		logs.Error("GetPackBody() got err=(%v)", err)
		return nil
	}
	return val
}

func (p *DataPack) SetPackBody(data []byte) {
	var err error
	p.data, err = p.bodyDecoder.Write(p.data, data)
	if err != nil {
		logs.Error("SetPackBody() got err=(%v)", err)
	}
}

// todo head + body
func (p *DataPack) String() string {
	cmd, ll, err := p.headDecoder.Read(p.data)
	if err != nil {
		logs.Error("GetCmdId() got err=(%v)", err)
		return ""
	}
	val, err := p.bodyDecoder.Read(p.data)
	if err != nil {
		logs.Error("GetCmdId() got err=(%v)", err)
		return ""
	}
	return fmt.Sprintf("cmd=%v, len=%v, body=%v", cmd, ll, string(val))
}
