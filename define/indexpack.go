/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: indexpack.go
 * @time: 2020/1/19 3:39 下午
 * @project: msgsubscribesvr
 */

package define

import (
	`fmt`

	`github.com/google/btree`

	`github.com/generalzgd/msg-subscriber/iface`
	`github.com/generalzgd/msg-subscriber/util`
)

type IndexPack struct {
	// 默认永远0
	Index uint64
	Data  uint64
}

func (p *IndexPack) String() string {
	return fmt.Sprintf(`{Index:%v, Data:%v}`, p.Index, p.Data)
}

func (p *IndexPack) SetRecordTime(t int64) {
}

func (p *IndexPack) GetCmdId() int {
	return 0
}

func (p *IndexPack) HasExpire(expire float64) bool {
	return false
}

func (p *IndexPack) GetPeers() []string {
	return nil
}

func (p *IndexPack) PutPeers(peers ...string) {
}

func (p *IndexPack) DelPeers(keys ...string) {
}

func (p *IndexPack) GetPeersRetry() map[string]uint32 {
	return nil
}

func (p *IndexPack) MergerPeers(item iface.StoreItem) {

}

func (p *IndexPack) Less(than btree.Item) bool {
	if p.Index < than.(*IndexPack).Index {
		return true
	}
	return false
}

func (p *IndexPack) GetIndex() uint64 {
	return p.Index
}

// 设置存储的key，仅限于空结构体使用
func (p *IndexPack) SetIndex(v uint64) {
	p.Index = v
}

func (p *IndexPack) GetRetry(peer string) uint32 {
	return 0
}

func (p *IndexPack) MaxRetry() uint32 {
	return 0
}

func (p *IndexPack) Serialize() []byte {
	return util.Uint64ToBytes(p.Data)
}

func (p *IndexPack) Deserialize(in []byte) error {
	p.Data = util.BytesToUint64(in)
	return nil
}
