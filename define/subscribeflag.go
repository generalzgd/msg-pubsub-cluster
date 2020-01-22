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

	`github.com/generalzgd/msg-subscriber/iproto`
)

// 标记集群内，目标协议是否有订阅过。采用数字cmdid作为key
type SubscribeFlagMap struct {
	lock sync.RWMutex
	data map[int]map[string]map[string]struct{} // cmdid => nodeId => consumerKye => nil
}

func NewSubscribeFlagMap() *SubscribeFlagMap {
	return &SubscribeFlagMap{
		data: map[int]map[string]map[string]struct{}{},
	}
}

func (p *SubscribeFlagMap) Copy() map[string]*iproto.NodeMap {
	p.lock.RLock()
	defer p.lock.RUnlock()

	out := make(map[string]*iproto.NodeMap, len(p.data))
	for id, nodes := range p.data {
		for nodeId, node := range nodes {
			nodeMap, ok := out[nodeId]
			if !ok {
				nodeMap = &iproto.NodeMap{
					Data: map[string]*iproto.IdList{},
				}
				out[nodeId] = nodeMap
			}
			for key := range node {
				keyMap, ok := nodeMap.Data[key]
				if !ok {
					keyMap = &iproto.IdList{
						Ids: []uint32{},
					}
					nodeMap.Data[key] = keyMap
				}
				keyMap.Ids = append(keyMap.Ids, uint32(id))
			}
		}
	}
	return out
}

func (p *SubscribeFlagMap) Set(nodesMap map[string]*iproto.NodeMap) {
	if len(nodesMap) < 1 {
		return
	}
	p.lock.Lock()
	defer p.lock.Unlock()

	for nodeId, nodeMap := range nodesMap {
		for key, list := range nodeMap.Data {
			for _, id := range list.Ids {
				nodes, ok := p.data[int(id)]
				if !ok {
					nodes = map[string]map[string]struct{}{}
					p.data[int(id)] = nodes
				}
				node, ok := nodes[nodeId]
				if !ok {
					node = map[string]struct{}{}
					nodes[nodeId] = node
				}
				node[key] = struct{}{}
			}
		}
	}
}

func (p *SubscribeFlagMap) Put(nodeId string, key string, ids []int) {
	if len(ids) < 1 {
		return
	}
	p.lock.Lock()
	defer p.lock.Unlock()

	for _, id := range ids {
		nodes, ok := p.data[id]
		if !ok {
			nodes = map[string]map[string]struct{}{}
			p.data[id] = nodes
		}
		node, ok := nodes[nodeId]
		if !ok {
			node = map[string]struct{}{}
			nodes[nodeId] = node
		}
		node[key] = struct{}{}
	}
}

func (p *SubscribeFlagMap) Del(nodeId string, key string, ids []int) {
	if len(ids) < 1 {
		return
	}
	p.lock.Lock()
	defer p.lock.Unlock()

	for _, id := range ids {
		nodes, ok := p.data[id]
		if !ok {
			continue
		}
		node, ok := nodes[nodeId]
		if !ok {
			continue
		}
		delete(node, key)

		if len(node) == 0 {
			delete(nodes, nodeId)
			if len(nodes) == 0 {
				delete(p.data, id)
			}
		}
	}
}

func (p *SubscribeFlagMap) Has(id int) bool {
	p.lock.RLock()
	defer p.lock.RUnlock()

	nodes, ok := p.data[id]
	if ok && len(nodes) > 0 {
		return true
	}
	return false
}

func (p *SubscribeFlagMap) GetSubscribedCount(id int) uint32 {
	p.lock.RLock()
	defer p.lock.RUnlock()
	nodes, ok := p.data[id]
	out := 0
	if ok {
		for _, node := range nodes {
			out += len(node)
		}
	}
	return uint32(out)
}
