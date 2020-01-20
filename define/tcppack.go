/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: tcppack.go
 * @time: 2020/1/20 11:13 上午
 * @project: msg-subscriber
 */

package define

import (
	`encoding/binary`
)

type PostPacket struct {
	Length   uint32 // body字节数
	ToType   uint16 // 转发id
	ToIp     uint32 //
	FromType uint16 //
	FromIp   uint32 //
	AskId    uint32 // 请求标记
	AskId2   uint16 //
	CmdId    uint16 // 协议id
	Body     []byte //
}

func (p *PostPacket) Serialize() []byte {
	buf := make([]byte, 24+len(p.Body))

	p.Length = uint32(len(p.Body))

	binary.LittleEndian.PutUint32(buf[0:], p.Length)
	binary.LittleEndian.PutUint16(buf[4:], p.ToType)
	binary.LittleEndian.PutUint32(buf[6:], p.ToIp)
	binary.LittleEndian.PutUint16(buf[10:], p.FromType)
	binary.LittleEndian.PutUint32(buf[12:], p.FromIp)
	binary.LittleEndian.PutUint32(buf[16:], p.AskId)
	binary.LittleEndian.PutUint16(buf[20:], p.AskId2)
	binary.LittleEndian.PutUint16(buf[22:], p.CmdId)

	copy(buf[24:], p.Body)
	return buf
}
