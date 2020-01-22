/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: decoder.go
 * @time: 2020/1/19 3:12 下午
 * @project: msgsubscribesvr
 */

package body

import (
	`encoding/binary`
	`errors`
)

type Decoder struct {
	pos   int
	order binary.ByteOrder
}

func NewDecoder(pos int, order binary.ByteOrder) *Decoder {
	return &Decoder{
		pos:   pos,
		order: order,
	}
}

func (p *Decoder) Read(data []byte) ([]byte, error) {
	if p.pos > len(data) {
		return nil, errors.New("cross the border")
	}
	return data[p.pos:], nil
}

func (p *Decoder) Write(buf, data []byte) ([]byte, error) {
	buf = append(buf[:p.pos], data...)
	return buf, nil
}