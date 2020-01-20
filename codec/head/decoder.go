/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: decoder.go
 * @time: 2020/1/19 3:08 下午
 * @project: msgsubscribesvr
 */

package head

import (
	`encoding/binary`
	`errors`
)

type Decoder struct {
	lenPos  int
	lenSize int
	cmdPos  int              // 起始读取的坐标
	cmdSize int              // 要读取字节的数量
	order   binary.ByteOrder // 大小端
	//
	cmdVal int
	lenVal int
}

func NewDecoder(cmdPos, cmdSize, lenPos, lenSize int, order binary.ByteOrder) Decoder {
	return Decoder{
		cmdPos:  cmdPos,
		cmdSize: cmdSize,
		lenPos:  lenPos,
		lenSize: lenSize,
		order:   order,
	}
}

func (p *Decoder) Read(data []byte) (int, int, error) {
	if p.cmdPos+p.cmdSize > len(data) || p.lenPos+p.lenSize > len(data) {
		return 0,0,errors.New("cross the border")
	}
	buf := data[p.cmdPos:]
	switch p.cmdSize {
	case 1:
		p.cmdVal = int(buf[0])
	case 2:
		p.cmdVal = int(p.order.Uint16(buf))
	case 4:
		p.cmdVal = int(p.order.Uint32(buf))
	case 8:
		p.cmdVal = int(p.order.Uint64(buf))
	}

	buf = data[p.lenPos:]
	switch p.lenSize {
	case 1:
		p.lenVal = int(buf[0])
	case 2:
		p.lenVal = int(p.order.Uint16(buf))
	case 4:
		p.lenVal = int(p.order.Uint32(buf))
	case 8:
		p.lenVal = int(p.order.Uint64(buf))
	}
	return p.cmdVal, p.lenVal, nil
}

func (p *Decoder) Write(buf []byte, cmdVal, lenVal int) error {
	if p.cmdPos+p.cmdSize > len(buf) || p.lenPos+p.lenSize > len(buf) {
		return errors.New("cross the border")
	}
	switch p.cmdSize {
	case 1:
		buf[p.cmdPos] = byte(cmdVal)
	case 2:
		p.order.PutUint16(buf[p.cmdPos:], uint16(cmdVal))
	case 4:
		p.order.PutUint32(buf[p.cmdPos:], uint32(cmdVal))
	case 8:
		p.order.PutUint64(buf[p.cmdPos:], uint64(cmdVal))
	}
	//
	switch p.lenSize {
	case 1:
		buf[p.lenPos] = byte(cmdVal)
	case 2:
		p.order.PutUint16(buf[p.lenPos:], uint16(lenVal))
	case 4:
		p.order.PutUint32(buf[p.lenPos:], uint32(lenVal))
	case 8:
		p.order.PutUint64(buf[p.lenPos:], uint64(lenVal))
	}
	return nil
}
