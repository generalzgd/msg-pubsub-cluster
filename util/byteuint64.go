/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: byteuint64.go
 * @time: 2020/1/2 2:25 下午
 * @project: msgsubscribesvr
 */

package util

import (
	`encoding/binary`
)

// Converts bytes to an integer
func BytesToUint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

// Converts a uint to a byte slice
func Uint64ToBytes(u uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, u)
	return buf
}


func GetIntFromBuf(buf []byte, pos, size int, order binary.ByteOrder) int {
	start := buf[pos:]
	val := 0
	switch size {
	case 1:
		val = int(start[0])
	case 2:
		val = int(order.Uint16(start))
	case 4:
		val = int(order.Uint32(start))
	case 8:
		val = int(order.Uint64(start))
	}
	return val
}

func SetIntToBuf(buf []byte, pos, size, val int, order binary.ByteOrder) {
	switch size {
	case 1:
		buf[pos] = byte(val)
	case 2:
		order.PutUint16(buf[pos:], uint16(val))
	case 4:
		order.PutUint32(buf[pos:], uint32(val))
	case 8:
		order.PutUint64(buf[pos:], uint64(val))
	}
}