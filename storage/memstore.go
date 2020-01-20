/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: MemBTree.go
 * @time: 2019/12/24 1:36 下午
 * @project: packagesubscribesvr
 */

package storage

import (
	`github.com/google/btree`

	`github.com/generalzgd/msg-subscriber/define`
	`github.com/generalzgd/msg-subscriber/iface`
)

type InMemStore struct {
	buckets     map[string]*btree.BTree
	itemFactory func(string) iface.StoreItem
}

func NewInMemStore(degree int, itemFactory func(string) iface.StoreItem, buckets ...string) *InMemStore {
	obj := &InMemStore{
		buckets:     map[string]*btree.BTree{},
		itemFactory: itemFactory,
	}
	for _, item := range buckets {
		obj.buckets[item] = btree.New(degree)
	}
	return obj
}

func (p *InMemStore) Close() error {
	return nil
}

func (p *InMemStore) Store(bucket string, val iface.StoreItem) error {
	if val == nil {
		return define.ParamNil
	}
	tree, ok := p.buckets[bucket]
	if !ok {
		return define.BucketEmpty
	}
	tree.ReplaceOrInsert(val)
	return nil
}

func (p *InMemStore) StoreBatch(bucket string, batch ...iface.StoreItem) error {
	tree, ok := p.buckets[bucket]
	if !ok {
		return define.BucketEmpty
	}
	for _, it := range batch {
		if it == nil {
			continue
		}
		tree.ReplaceOrInsert(it)
	}
	return nil
}

func (p *InMemStore) UpdateBatch(bucket string, batch ...iface.StoreItem) error {
	tree, ok := p.buckets[bucket]
	if !ok {
		return define.BucketEmpty
	}
	for _, it := range batch {
		if it == nil {
			continue
		}
		tp := tree.Get(it)
		if tp == nil {
			tree.ReplaceOrInsert(it)
			continue
		}
		obj, ok := tp.(iface.StoreItem)
		if !ok {
			continue
		}
		it.MergerPeers(obj)
		tree.ReplaceOrInsert(it)
	}
	return nil
}

// 删除[min, max)之间的值
func (p *InMemStore) DeleteRange(bucket string, min, max iface.StoreItem) error {
	if min == nil || max == nil {
		return define.ParamNil
	}
	tree, ok := p.buckets[bucket]
	if !ok {
		return define.BucketEmpty
	}
	list := make([]iface.StoreItem, 0, max.GetIndex()-min.GetIndex())
	tree.AscendRange(min, max, func(i btree.Item) bool {
		list = append(list, i.(iface.StoreItem))
		return true
	})

	for _, it := range list {
		tree.Delete(it)
	}
	return nil
}

func (p *InMemStore) Delete(bucket string, items ...iface.StoreItem) error {
	tree, ok := p.buckets[bucket]
	if !ok {
		return define.BucketEmpty
	}
	for _, it := range items {
		if it == nil {
			continue
		}
		tree.Delete(it)
	}
	return nil
}

func (p *InMemStore) DeleteByIndexes(bucket string, indexes ...uint64) ([]iface.StoreItem, error) {
	tree, ok := p.buckets[bucket]
	if !ok {
		return nil, define.BucketEmpty
	}

	out := make([]iface.StoreItem, 0, len(indexes))
	for _, index := range indexes {
		key := p.itemFactory(bucket)
		key.SetIndex(index)
		val := tree.Get(key)
		if val != nil {
			out = append(out, val.(iface.StoreItem))
		}
	}
	for _, it := range out {
		tree.Delete(it)
	}
	return out, nil
}

func (p *InMemStore) GetUint64(bucket, key string) (uint64, error) {
	tree, ok := p.buckets[bucket]
	if !ok {
		return 0, define.BucketEmpty
	}
	it := tree.Get(&define.IndexPack{})
	if it != nil {
		res, ok := it.(iface.StoreItem)
		if ok {
			return res.GetIndex(), nil
		}
	}
	return 0, nil
}

func (p *InMemStore) SetUnit64(bucket, key string, val uint64) (uint64, error) {
	tree, ok := p.buckets[bucket]
	if !ok {
		return 0, define.BucketEmpty
	}
	it := &define.IndexPack{
		Data: val,
	}
	tree.ReplaceOrInsert(it)
	return val, nil
}

func (p *InMemStore) GetBatch(bucket string, limit int) ([]iface.StoreItem, error) {
	tree, ok := p.buckets[bucket]
	if !ok {
		return nil, define.BucketEmpty
	}
	ll, min := tree.Len(), tree.Min()
	if limit > ll {
		limit = ll
	}
	if limit < 1 {
		return nil, nil
	}

	out := make([]iface.StoreItem, 0, limit)
	cnt := 0
	tree.AscendGreaterOrEqual(min, func(i btree.Item) bool {
		out = append(out, i.(iface.StoreItem))
		cnt++
		if cnt <= limit {
			return true
		}
		return false
	})

	// for _, it := range out {
	// 	tree.Delete(it)
	// }
	return out, nil
}

func (p *InMemStore) GetBatchByIndexes(bucket string, indexes ...uint64) ([]iface.StoreItem, error) {
	tree, ok := p.buckets[bucket]
	if !ok {
		return nil, define.BucketEmpty
	}

	out := make([]iface.StoreItem, 0, len(indexes))
	for _, index := range indexes {
		key := p.itemFactory(bucket)
		key.SetIndex(index)
		val := tree.Get(key)
		if val != nil {
			out = append(out, val.(iface.StoreItem))
		}
	}
	return out, nil
}

func (p *InMemStore) GetMin(bucket string) iface.StoreItem {
	tree, ok := p.buckets[bucket]
	if !ok {
		return nil
	}
	v := tree.Min()
	return v.(iface.StoreItem)
}

func (p *InMemStore) GetMax(bucket string) iface.StoreItem {
	tree, ok := p.buckets[bucket]
	if !ok {
		return nil
	}
	v := tree.Max()
	return v.(iface.StoreItem)
}
