/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: decoder.go
 * @time: 2020/1/19 3:42 下午
 * @project: msgsubscribesvr
 */

package cmdstr

import (
	`bytes`
	`encoding/json`
)

type tmpCmd struct {
	CmdId string `json:"cmdid"`
}


type Decoder struct {
}

func NewDecoder() Decoder {
	return Decoder{}
}

func (p *Decoder) Read(data []byte) (string, error) {
	args := tmpCmd{}
	err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&args)
	if err != nil {
		return "", err
	}
	return args.CmdId, nil
}
