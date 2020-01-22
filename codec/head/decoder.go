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

	`github.com/generalzgd/msg-subscriber/util`
)

type Decoder struct {
	cmdPos  int // 起始读取的坐标
	cmdSize int // 要读取字节的数量
	lenPos  int
	lenSize int
	order   binary.ByteOrder // 大小端
	//
	cmd int
	len int
}

func NewDecoder(cmdPos, cmdSize, lenPos, lenSize int, order binary.ByteOrder) *Decoder {
	return &Decoder{
		cmdPos:  cmdPos,
		cmdSize: cmdSize,
		lenPos:  lenPos,
		lenSize: lenSize,
		order:   order,
	}
}

func (p *Decoder) Read(data []byte) (int, int, error) {
	if p.cmdPos+p.cmdSize > len(data) || p.lenPos+p.lenSize > len(data) {
		return 0,0, errors.New("cross the border")
	}
	p.cmd = util.GetIntFromBuf(data, p.cmdPos, p.cmdSize, p.order)
	p.len = util.GetIntFromBuf(data, p.lenPos, p.lenSize, p.order)
	return p.cmd, p.len, nil
}

func (p *Decoder) Write(buf []byte, cmd, l int) error {
	if p.cmdPos+p.cmdSize > len(buf) || p.lenPos+p.lenSize > len(buf) {
		return errors.New("cross the border")
	}
	util.SetIntToBuf(buf, p.cmdPos, p.cmdSize, cmd, p.order)
	util.SetIntToBuf(buf, p.lenPos, p.lenSize, l, p.order)
	return nil
}
