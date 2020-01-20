/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: boltstore.go
 * @time: 2019/12/24 5:26 下午
 * @project: packagesubscribesvr
 */

package storage

import (
	`github.com/astaxie/beego/logs`
	`github.com/boltdb/bolt`

	`github.com/generalzgd/msg-subscriber/define`
	`github.com/generalzgd/msg-subscriber/iface`
	`github.com/generalzgd/msg-subscriber/util`
)

type BoltStore struct {
	conn        *bolt.DB
	path        string
	itemFactory func(string) iface.StoreItem
}

func NewBoltStore(path string, itemFactory func(string) iface.StoreItem, buckets ...string) *BoltStore {
	// Try to connect
	handle, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil
	}

	// Create the new store
	store := &BoltStore{
		conn:        handle,
		path:        path,
		itemFactory: itemFactory,
	}
	if err := store.initialize(buckets...); err != nil {
		store.Close()
		return nil
	}
	return store
}

func (p *BoltStore) Close() error {
	return p.conn.Close()
}

func (p *BoltStore) Store(bucket string, val iface.StoreItem) error {
	if val == nil {
		return define.ParamNil
	}
	tx, err := p.conn.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	bk := tx.Bucket([]byte(bucket))
	if err := bk.Put(util.Uint64ToBytes(val.GetIndex()), val.Serialize()); err != nil {
		return err
	}

	return tx.Commit()
}

func (p *BoltStore) StoreBatch(bucket string, batch ...iface.StoreItem) error {
	tx, err := p.conn.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	bk := tx.Bucket([]byte(bucket))
	for _, val := range batch {
		if val == nil {
			continue
		}
		if err := bk.Put(util.Uint64ToBytes(val.GetIndex()), val.Serialize()); err != nil {
			logs.Error("bolt store save error.", val)
		}
	}

	return tx.Commit()
}

func (p *BoltStore) UpdateBatch(bucket string, batch ...iface.StoreItem) error {
	tx, err := p.conn.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	bk := tx.Bucket([]byte(bucket))

	for _, it := range batch {
		if it == nil {
			continue
		}
		val := bk.Get(util.Uint64ToBytes(it.GetIndex()))
		if val == nil {
			bk.Put(util.Uint64ToBytes(it.GetIndex()), it.Serialize())
			continue
		}
		item := p.itemFactory(bucket)
		if err := item.Deserialize(val); err != nil {
			continue
		}
		it.MergerPeers(item)
		bk.Put(util.Uint64ToBytes(it.GetIndex()), it.Serialize())
	}
	return tx.Commit()
}

// initialize is used to set up all of the buckets.
func (p *BoltStore) initialize(buckets ...string) error {
	tx, err := p.conn.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	// Create all the buckets
	for _, bucket := range buckets {
		if _, err := tx.CreateBucketIfNotExists([]byte(bucket)); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// FirstIndex returns the first known index from the Raft log.
func (p *BoltStore) firstIndex(bucket string) (string, error) {
	tx, err := p.conn.Begin(false)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	curs := tx.Bucket([]byte(bucket)).Cursor()
	if first, _ := curs.First(); first == nil {
		return "", nil
	} else {
		return string(first), nil
	}
}

// LastIndex returns the last known index from the Raft log.
func (p *BoltStore) lastIndex(bucket string) (string, error) {
	tx, err := p.conn.Begin(false)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	curs := tx.Bucket([]byte(bucket)).Cursor()
	if last, _ := curs.Last(); last == nil {
		return "", nil
	} else {
		return string(last), nil
	}
}

func (p *BoltStore) Get(bucket, key string) ([]byte, error) {
	tx, err := p.conn.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	val := tx.Bucket([]byte(bucket)).Get([]byte(key))
	return val, nil
}

func (p *BoltStore) GetUint64(bucket, key string) (uint64, error) {
	tx, err := p.conn.Begin(false)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	val := tx.Bucket([]byte(bucket)).Get([]byte(key))
	return util.BytesToUint64(val), nil
}

func (p *BoltStore) SetUnit64(bucket, key string, val uint64) (uint64, error) {
	tx, err := p.conn.Begin(true)
	if err != nil {
		return val, err
	}
	defer tx.Rollback()

	err = tx.Bucket([]byte(bucket)).Put([]byte(key), util.Uint64ToBytes(val))
	if err != nil {
		return val, err
	}
	return val, tx.Commit()
}

// Get is used to retrieve a value from the k/v store by key
func (p *BoltStore) GetBatch(bucket string, limit int) ([]iface.StoreItem, error) {
	tx, err := p.conn.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	curs := tx.Bucket([]byte(bucket)).Cursor()
	out := make([]iface.StoreItem, 0, 100)

	for k, v := curs.First(); k != nil; k, v = curs.Next() {

		item := p.itemFactory(bucket)
		if err := item.Deserialize(v); err == nil {
			out = append(out, item)
			continue
		}
		// todo 处理成功才能删除
		// if err := curs.Delete(); err != nil {
		// 	continue
		// }
	}

	return out, nil
}

func (p *BoltStore) GetBatchByIndexes(bucket string, indexes ...uint64) ([]iface.StoreItem, error) {
	tx, err := p.conn.Begin(false)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	bk := tx.Bucket([]byte(bucket))

	out := make([]iface.StoreItem,0,len(indexes))

	for _, index := range indexes {
		val := bk.Get(util.Uint64ToBytes(index))
		if val == nil {
			continue
		}
		item := p.itemFactory(bucket)
		if err := item.Deserialize(val); err != nil {
			continue
		}
		out = append(out, item)
	}
	return out, nil
}

func (p *BoltStore) DeleteRange(bucket string, min, max iface.StoreItem) error {
	if min == nil || max == nil {
		return define.ParamNil
	}
	tx, err := p.conn.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	bk := tx.Bucket([]byte(bucket))
	curs := bk.Cursor()
	for k, _ := curs.Seek(util.Uint64ToBytes(min.GetIndex())); k != nil; k, _ = curs.Next() {
		// Handle out-of-range log index
		if util.BytesToUint64(k) > max.GetIndex() {
			break
		}

		// Delete in-range log index
		if err := curs.Delete(); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (p *BoltStore) Delete(bucket string, items ...iface.StoreItem) error {
	tx, err := p.conn.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	bk := tx.Bucket([]byte(bucket))
	for _, it := range items {
		if it == nil {
			continue
		}
		bk.Delete(util.Uint64ToBytes(it.GetIndex()))
	}

	return tx.Commit()
}

func (p *BoltStore) DeleteByIndexes(bucket string, indexes ...uint64) ([]iface.StoreItem, error) {
	tx, err := p.conn.Begin(true)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	bk := tx.Bucket([]byte(bucket))

	out := make([]iface.StoreItem,0,len(indexes))

	for _, index := range indexes {
		key := util.Uint64ToBytes(index)
		val := bk.Get(key)
		if val == nil {
			continue
		}
		item := p.itemFactory(bucket)
		if err := item.Deserialize(val); err != nil {
			continue
		}
		out = append(out, item)
		//
		bk.Delete(key)
	}
	return out, nil
}

// 获取最小的值
func (p *BoltStore) GetMin(bucket string) iface.StoreItem {
	tx, err := p.conn.Begin(false)
	if err != nil {
		return nil
	}
	defer tx.Rollback()

	curs := tx.Bucket([]byte(bucket)).Cursor()
	first, bts := curs.First()
	if first == nil {
		return nil
	}
	item := p.itemFactory(bucket)
	if err := item.Deserialize(bts); err != nil {
		return nil
	}
	return item
}

func (p *BoltStore) GetMax(bucket string) iface.StoreItem {
	tx, err := p.conn.Begin(false)
	if err != nil {
		return nil
	}
	defer tx.Rollback()

	curs := tx.Bucket([]byte(bucket)).Cursor()
	last, bts := curs.Last()
	if last == nil {
		return nil
	}
	item := p.itemFactory(bucket)
	if err := item.Deserialize(bts); err != nil {
		return nil
	}
	return item
}
