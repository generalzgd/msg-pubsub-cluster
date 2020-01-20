/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: subscribeflag.go
 * @time: 2020/1/2 11:30 上午
 * @project: msgsubscribesvr
 */

package define

import (
	`sync`
)

// 标记集群内，目标协议是否有订阅过。采用数字cmdid作为key
type SubscribeFlagMap struct {
	lock sync.RWMutex
	data map[int]uint32 // cmdid =>
}

func NewSubscribeFlagMap() *SubscribeFlagMap {
	return &SubscribeFlagMap{
		data: map[int]uint32{},
	}
}

func (p *SubscribeFlagMap) Copy() map[uint32]uint32 {
	p.lock.RLock()
	defer p.lock.RUnlock()

	out := make(map[uint32]uint32, len(p.data))
	for k, v := range p.data {
		out[uint32(k)] = v
	}
	return out
}

func (p *SubscribeFlagMap) Set(m map[uint32]uint32) {
	if len(m) < 1 {
		return
	}
	p.lock.Lock()
	defer p.lock.Unlock()

	for k, v := range m {
		pv, ok := p.data[int(k)]
		if ok {
			if v > pv {
				p.data[int(k)] = v
			}
		} else {
			p.data[int(k)] = v
		}
	}
}

func (p *SubscribeFlagMap) Put(cmdIds ...int) {
	if len(cmdIds) < 1 {
		return
	}
	p.lock.Lock()
	defer p.lock.Unlock()

	for _, id := range cmdIds {
		v, ok := p.data[id]
		if ok {
			p.data[id] = v + 1
		} else {
			p.data[id] = 1
		}
		//p.data.Store(id, struct {}{}) // 要改成计数
	}
}

func (p *SubscribeFlagMap) Del(cmdIds ...int) {
	if len(cmdIds) < 1 {
		return
	}
	p.lock.Lock()
	defer p.lock.Unlock()

	for _, id := range cmdIds {
		v, ok := p.data[id]
		if ok && v > 0 {
			if v > 1 {
				p.data[id] = v - 1
			} else {
				delete(p.data, id)
			}
		} else {
			delete(p.data, id)
		}
	}
}

func (p *SubscribeFlagMap) Has(id int) bool {
	p.lock.RLock()
	defer p.lock.RUnlock()

	v, ok := p.data[id]
	if ok && v > 0 {
		return true
	}
	return false
}

func (p *SubscribeFlagMap) GetSubscribedCount(id int) uint32 {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.data[id]
}
