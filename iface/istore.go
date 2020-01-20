/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: istore.go
 * @time: 2020/1/2 2:08 下午
 * @project: msgsubscribesvr
 */

package iface

import (
	`fmt`
	`io`

	`github.com/google/btree`
)

// 统一的数据流接口  -->  1. memStore
// 						2. boltStore
// store：
// 1. 实现排序功能
// 2. memory速度快
// 3. bolt实现文件存储，以防数据丢失

type StoreItem interface {
	fmt.Stringer
	btree.Item
	SetRecordTime(int64)
	GetCmdId() int
	HasExpire(float64) bool
	GetPeers() []string
	PutPeers(...string)
	DelPeers(...string)
	GetPeersRetry() map[string]uint32
	MergerPeers(StoreItem)
	GetIndex() uint64
	SetIndex(v uint64)
	GetRetry(string) uint32
	MaxRetry() uint32
	Serialize() []byte
	Deserialize([]byte) error
}

// 存储桥
type IStoreBridge interface {
	io.Closer
	Store(StoreItem) error
	StoreBatch(...StoreItem) error
	UpdateBatch(...StoreItem) error
	GetBatch(int) ([]StoreItem, error)
	GetBatchByIndexes(...uint64) ([]StoreItem, error)
	DeleteRange(min, max StoreItem) error
	Delete(...StoreItem) error
	DeleteByIndexes(indexes ...uint64) ([]StoreItem, error)
	GetUint64(string) (uint64, error)
	SetUnit64(string, uint64) (uint64, error)
}

type IStoreOperator interface {
	io.Closer
	Store(string, StoreItem) error
	StoreBatch(string, ...StoreItem) error
	UpdateBatch(string, ...StoreItem) error
	DeleteRange(bucket string, min, max StoreItem) error
	Delete(string, ...StoreItem) error
	DeleteByIndexes(bucket string, indexes ...uint64) ([]StoreItem, error)
	GetBatch(string, int) ([]StoreItem, error)
	GetBatchByIndexes(bucket string, indexes ...uint64) ([]StoreItem, error)
	GetUint64(bucket, key string) (uint64, error)
	SetUnit64(bucket, key string, val uint64) (uint64, error)
}
