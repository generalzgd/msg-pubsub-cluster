/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: subscribemao.go
 * @time: 2019/12/24 10:04 上午
 * @project: packagesubscribesvr
 */

package define

import (
	`sync`

	`github.com/generalzgd/msg-subscriber/iface`
)

// ********************************
// 当前节点的订阅信息，单条信息可订阅多个消费者
type SubscribeMap struct {
	lock sync.RWMutex
	// 协议ID转consumer
	id2Consumer map[int][]iface.IConsumer
	//
	// consumer key转协议ID
	consumer2Ids map[string]map[int]struct{} // 节点对应协议id, {NodeKey: CmdId:struct}
	// consumer key 转consumer
	key2Consumer map[string]iface.IConsumer
}

func NewSubscribeMap() *SubscribeMap {
	return &SubscribeMap{
		//cmdMap:   map[string][]*SubscribeInfo{},
		id2Consumer:  map[int][]iface.IConsumer{},
		consumer2Ids: map[string]map[int]struct{}{},
		key2Consumer: map[string]iface.IConsumer{},
	}
}

// 只要找到一个包含
func (p *SubscribeMap) has(list []iface.IConsumer, tar iface.IConsumer) bool {
	for _, it := range list {
		if it.Equal(tar) {
			return true
		}
	}
	return false
}

// 只删除找到的第一个
func (p *SubscribeMap) remove(list []iface.IConsumer, tar iface.IConsumer) ([]iface.IConsumer, bool) {
	for i, it := range list {
		if it.Equal(tar) {
			// 删除最后一个
			if i == len(list)-1 {
				return list[:i], true
			}
			// 删除中间的
			return append(list[:i], list[i+1]), true
		}
	}
	return list, false
}

func (p *SubscribeMap) Put(ids []int, v iface.IConsumer) []int {
	p.lock.Lock()
	defer p.lock.Unlock()
	putted := make([]int, 0, len(ids))

	for _, id := range ids {
		m, ok := p.id2Consumer[id]
		if ok {
			if !p.has(m, v) {
				m = append(m, v)
				p.id2Consumer[id] = m
				putted = append(putted, id)
			}
		} else {
			p.id2Consumer[id] = []iface.IConsumer{v}
			putted = append(putted, id)
		}
	}
	//
	if len(putted) > 0 {
		k := v.GetKey()
		m, ok := p.consumer2Ids[k]
		if !ok {
			m = map[int]struct{}{}
		}
		for _, id := range putted {
			m[id] = struct{}{}
		}
		p.consumer2Ids[k] = m
		//
		if t, ok := p.key2Consumer[k]; !ok || !t.Equal(v) {
			p.key2Consumer[k] = v
		}
	}
	return putted
}

func (p *SubscribeMap) Del(ids []int, v iface.IConsumer) []int {
	p.lock.Lock()
	defer p.lock.Unlock()

	removed := make([]int, 0, len(ids))
	for _, id := range ids {
		if m, ok := p.id2Consumer[id]; ok {
			if ok {
				if l, ok := p.remove(m, v); ok {
					if len(l) > 0 {
						p.id2Consumer[id] = l
					} else {
						delete(p.id2Consumer, id)
					}
					removed = append(removed, id)
				}
			}
		}
	}
	//
	if len(removed) > 0 {
		k := v.GetKey()
		m, ok := p.consumer2Ids[k]
		if ok {
			for _, id := range removed {
				delete(m, id)
			}
			if len(m) == 0 {
				// 该节点已经没有注册协议，
				delete(p.consumer2Ids, k)
				// 该节点已经没有注册协议
				delete(p.key2Consumer, k)
			} else {
				p.consumer2Ids[k] = m
			}
		}
	}
	return removed
}

// 通过注册的协议ID，获取对应的订阅者（同一个消息可能有不同的订阅者）
func (p *SubscribeMap) GetConsumer(id int) ([]iface.IConsumer, bool) {
	p.lock.RLock()
	defer p.lock.RUnlock()
	m, ok := p.id2Consumer[id]
	if ok {
		return m, ok
	}
	return nil, false
}

// 通过消费者的唯一可以获取对应的订阅者信息
func (p *SubscribeMap) GetConsumerByKey(keys ...string) []iface.IConsumer {
	p.lock.RLock()
	defer p.lock.RUnlock()

	out := make([]iface.IConsumer, 0, len(keys))
	for _, key := range keys {
		if v, ok := p.key2Consumer[key]; ok {
			out = append(out, v)
		}
	}
	return out
}

// 检测对应的消息ID是否有订阅者订阅
func (p *SubscribeMap) HasBooked(id int) bool {
	p.lock.RLock()
	defer p.lock.RUnlock()

	m, ok := p.id2Consumer[id]
	if ok && len(m) > 0 {
		return true
	}
	return false
}

// 获取订阅者订阅的所有协议ID
func (p *SubscribeMap) GetIdsByConsumer(tar iface.IConsumer) []int {
	p.lock.RLock()
	defer p.lock.RUnlock()

	m, ok := p.consumer2Ids[tar.GetKey()]
	if ok {
		out := make([]int, 0, len(m))
		for k := range m {
			out = append(out, k)
		}
		return out
	}
	return nil
}
