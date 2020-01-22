/**
 * @version: 1.0.0
 * @author: zhangguodong:general_zgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: iheadcodec.go
 * @time: 2020/1/21 11:26
 */
package codec

type IHeadCodec interface {
	Read(data []byte) (int, int, error)
	Write(buf []byte, cmd, l int) error
}
