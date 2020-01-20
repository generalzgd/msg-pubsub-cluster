/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: iflowstore.go
 * @time: 2019/12/24 5:40 下午
 * @project: packagesubscribesvr
 */

package storage

import (

	`github.com/generalzgd/msg-subscriber/define`
	`github.com/generalzgd/msg-subscriber/iface`
)

const (
	// 内存模型队列
	ModeInMem = 0
	// 文件持久化模型队列
	ModeFile = 1
)


func NewStoreOperator(mode int, degree int, path string, itemFactory func(string) iface.StoreItem, buckets ...string) iface.IStoreOperator {
	if mode == ModeInMem {
		return NewInMemStore(degree, itemFactory, buckets...)
	} else if mode == ModeFile {
		return NewBoltStore(path, itemFactory, buckets...)
	}
	return nil
}

// ***************************************************
type FlowStoreBridge struct {
	// 队列名字, 一个队列对应一个 bucket name
	bucketName string
	//
	store iface.IStoreOperator
}

func NewStoreBridge(bucket string, store iface.IStoreOperator) iface.IStoreBridge {
	return &FlowStoreBridge{
		bucketName: bucket,
		store:      store,
	}
}

func (p *FlowStoreBridge) Close() error {
	if p.store != nil {
		return p.store.Close()
	}
	return nil
}

func (p *FlowStoreBridge) Store(data iface.StoreItem) error {
	if data == nil {
		return nil
	}
	if p.store == nil {
		return define.StorageEmpty
	}
	return p.store.Store(p.bucketName, data)
}

func (p *FlowStoreBridge) StoreBatch(batch ...iface.StoreItem) error {
	if len(batch) == 0 {
		return nil
	}
	if p.store == nil {
		return define.StorageEmpty
	}
	return p.store.StoreBatch(p.bucketName, batch...)
}

func (p *FlowStoreBridge) UpdateBatch(batch ...iface.StoreItem) error {
	if len(batch) == 0 {
		return nil
	}
	if p.store == nil {
		return define.StorageEmpty
	}
	return p.store.UpdateBatch(p.bucketName, batch...)
}


func (p *FlowStoreBridge) GetBatch(limit int) ([]iface.StoreItem, error) {
	if p.store == nil {
		return nil, define.StorageEmpty
	}
	if limit < 1 {
		limit = 100
	}
	return p.store.GetBatch(p.bucketName, limit)
}

func (p *FlowStoreBridge) GetBatchByIndexes(indexes ...uint64) ([]iface.StoreItem, error) {
	if len(indexes) < 1 {
		return nil, nil
	}
	if p.store == nil {
		return nil, define.StorageEmpty
	}
	return p.store.GetBatchByIndexes(p.bucketName, indexes...)
}

func (p *FlowStoreBridge) DeleteRange(min, max iface.StoreItem) error {
	if min == nil || max == nil {
		return nil
	}
	if p.store == nil {
		return define.StorageEmpty
	}
	return p.store.DeleteRange(p.bucketName, min, max)
}

func (p *FlowStoreBridge) Delete(items ...iface.StoreItem) error {
	if len(items) == 0 {
		return nil
	}
	if p.store == nil {
		return define.StorageEmpty
	}
	return p.store.Delete(p.bucketName, items...)
}

func (p *FlowStoreBridge) DeleteByIndexes(indexes ...uint64) ([]iface.StoreItem, error) {
	if len(indexes) == 0 {
		return nil, nil
	}
	if p.store == nil {
		return nil, define.StorageEmpty
	}
	return p.store.DeleteByIndexes(p.bucketName, indexes...)
}

func (p *FlowStoreBridge) GetUint64(key string) (uint64, error) {
	if p.store == nil {
		return 0, define.StorageEmpty
	}
	return p.store.GetUint64(p.bucketName, key)
}

func (p *FlowStoreBridge) SetUnit64(key string, val uint64) (uint64, error) {
	if p.store == nil {
		return val, define.StorageEmpty
	}
	return p.store.SetUnit64(p.bucketName, key, val)
}
