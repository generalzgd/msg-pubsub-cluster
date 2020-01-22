/**
 * @version: 1.0.0
 * @author: zhangguodong:general_zgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: ibodycodec.go
 * @time: 2020/1/21 11:28
 */
package codec

type IBodyCodec interface {
	Read(data []byte) ([]byte, error)
	Write(buf, data []byte) ([]byte, error)
}
