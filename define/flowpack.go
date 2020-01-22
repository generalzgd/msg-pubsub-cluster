/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: def.go
 * @time: 2019/12/23 5:23 下午
 * @project: packagesubscribesvr
 */

package define

import (
	`fmt`
	`sync`
	`time`

	`github.com/golang/protobuf/proto`
	`github.com/google/btree`

	`github.com/generalzgd/msg-subscriber/codec`
	`github.com/generalzgd/msg-subscriber/codec/cmdstr`
	`github.com/generalzgd/msg-subscriber/iface`
	`github.com/generalzgd/msg-subscriber/iproto`
)

//
type FlowPack struct {
	Index      uint64
	peersLock  sync.RWMutex
	Peers      map[string]uint32 // 如果为空，表示所有的订阅消费者；存在则表示对应的消费者需要重发, 和已发几次
	Pack       *codec.DataPack
	RecordTime int64 // 记录时间戳
	//
	headDecoder codec.IHeadCodec
	bodyDecoder codec.IBodyCodec
}

func NewFlowPack(headDecoder codec.IHeadCodec, bodyDecoder codec.IBodyCodec) *FlowPack {
	return &FlowPack{
		Peers:       map[string]uint32{},
		headDecoder: headDecoder,
		bodyDecoder: bodyDecoder,
	}
}

func (p *FlowPack) String() string {
	return fmt.Sprintf(`{Index:%v, Data:%v}\n`, p.Index, p.Pack.String())
}

func (p *FlowPack) SetRecordTime(t int64) {
	// 死信阶段
	if len(p.Peers) > 0 {
		return
	}
	p.RecordTime = t
}

func (p *FlowPack) GetCmdId() int {
	return p.Pack.GetCmdId()
	//return v
}

// 是否过期
func (p *FlowPack) HasExpire(expire float64) bool {
	if p.RecordTime < 1 {
		return false
	}
	return float64(time.Now().Unix()) > expire+float64(p.RecordTime)
}

func (p *FlowPack) GetPeers() []string {
	p.peersLock.RLock()
	defer p.peersLock.RUnlock()

	out := make([]string, 0, len(p.Peers))
	for k := range p.Peers {
		out = append(out, k)
	}
	return out
}

// key为消费者订阅的时候提供的唯一key
func (p *FlowPack) PutPeers(keys ...string) {
	p.peersLock.Lock()
	defer p.peersLock.Unlock()

	if p.Peers == nil {
		p.Peers = map[string]uint32{}
	}

	for _, key := range keys {
		if v, ok := p.Peers[key]; ok {
			p.Peers[key] = v + 1
		} else {
			p.Peers[key] = 1
		}
	}
}

func (p *FlowPack) DelPeers(keys ...string) {
	p.peersLock.Lock()
	defer p.peersLock.Unlock()

	for _, key := range keys {
		delete(p.Peers, key)
	}
}

//
func (p *FlowPack) GetPeersRetry() map[string]uint32 {
	p.peersLock.RLock()
	defer p.peersLock.RUnlock()

	out := make(map[string]uint32, len(p.Peers))
	for k, v := range p.Peers {
		out[k] = v
	}
	return out
}

func (p *FlowPack) MergerPeers(item iface.StoreItem) {
	p.peersLock.Lock()
	defer p.peersLock.Unlock()

	tar := item.GetPeersRetry()
	// 如果自己没有，则保存目标值
	// 如果自己有，自己值小于目标值，则更新为目标值
	// 其他情况忽略
	for tk, tv := range tar {
		if v, ok := p.Peers[tk]; ok {
			if tv > v {
				p.Peers[tk] = tv
			}
		} else {
			p.Peers[tk] = tv
		}
	}
}

// 获取存储的key
func (p *FlowPack) GetIndex() uint64 {
	return p.Index
}

// 设置存储的key，仅限于空结构体使用
func (p *FlowPack) SetIndex(v uint64) {
	p.Index = v
}

// 获取重试次数
func (p *FlowPack) GetRetry(peer string) uint32 {
	if v, ok := p.Peers[peer]; ok {
		return v
	}
	return 0
}

func (p *FlowPack) MaxRetry() uint32 {
	max := uint32(0)
	for _, v := range p.Peers {
		if v > max {
			max = v
		}
	}
	return max
}

// 转换成StorePack字节
func (p *FlowPack) Serialize() []byte {
	//panic("implement me")
	pk := p.ToStorePack()
	out, err := proto.Marshal(&pk)
	if err != nil {
		return nil
	}
	//p.Data = out
	return out
}

// StorePack字节转换成flowpack
func (p *FlowPack) Deserialize(in []byte) error {
	//panic("implement me")
	//p.Data = in

	pk := iproto.StorePack{}
	if err := proto.Unmarshal(in, &pk); err != nil {
		return err
	}
	p.FromStorePack(pk)
	return nil
}

func (p *FlowPack) Less(than btree.Item) bool {
	if p.Index < than.(*FlowPack).Index {
		return true
	}
	return false
}

func (p *FlowPack) ToStorePack() iproto.StorePack {
	return iproto.StorePack{
		Index:      p.Index,
		Peers:      p.Peers,
		Data:       p.Pack.GetData(),
		RecordTime: p.RecordTime,
	}
}

func (p *FlowPack) FromStorePack(from iproto.StorePack) {
	p.Index = from.Index
	p.Peers = from.Peers
	//p.Data = from.Data
	p.Pack = codec.NewDataPack(from.Data, p.headDecoder, p.bodyDecoder)
	p.RecordTime = from.RecordTime
}

func (p *FlowPack) GetCmdStr() string {
	decoder := cmdstr.NewDecoder()
	val, err := decoder.Read(p.Pack.GetPackBody())
	if err != nil {
		return ""
	}
	return val
}
