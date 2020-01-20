/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: memstore_test.go.go
 * @time: 2019/12/26 6:21 下午
 * @project: msgsubscribesvr
 */

package storage

import (
	`testing`

	`github.com/generalzgd/msg-subscriber/define`
	`github.com/generalzgd/msg-subscriber/iface`
)

var (
	indexBucket = "index"
	saveBucket  = "test"
	memStore *InMemStore
)

func init() {
	memStore = NewInMemStore(3, func(s string) iface.StoreItem {
		if s == saveBucket {
			return &define.FlowPack{}
		}
		return &define.IndexPack{}
	}, indexBucket, saveBucket)
}

func TestInMemStore_Set_Get_Uint64(t *testing.T) {

	if _, err := memStore.SetUnit64(indexBucket, "", 1234); err != nil {
		t.Log("set error.", err)
		return
	}
	if got, err := memStore.GetUint64(indexBucket, ""); err != nil {
		t.Log("get error.", err)
		return
	} else {
		t.Log("got result.", got)
	}
}

func TestInMemStore_Store(t *testing.T) {
	type args struct {
		bucket string
		val    iface.StoreItem
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "TestInMemStore_Store",
			args: args{
				bucket: saveBucket,
				val: &define.FlowPack{
					Index: 11,
					Peers: map[string]uint32{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := memStore.Store(tt.args.bucket, tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("Store() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInMemStore_StoreBatch(t *testing.T) {
	type args struct {
		bucket string
		batch  []iface.StoreItem
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "TestInMemStore_StoreBatch",
			args: args{
				bucket: saveBucket,
				batch: []iface.StoreItem{
					&define.FlowPack{
						Index: 100,
						Peers: map[string]uint32{},
					},
					&define.FlowPack{
						Index: 101,
						Peers: map[string]uint32{},
					},
					&define.FlowPack{
						Index: 102,
						Peers: map[string]uint32{},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := memStore.StoreBatch(tt.args.bucket, tt.args.batch...); (err != nil) != tt.wantErr {
				t.Errorf("StoreBatch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInMemStore_UpdateBatch(t *testing.T) {
	type args struct {
		bucket string
		batch  []iface.StoreItem
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "TestInMemStore_UpdateBatch",
			args: args{
				bucket: saveBucket,
				batch: []iface.StoreItem{
					&define.FlowPack{
						Index: 1,
						Peers: map[string]uint32{
							"booker1": 1,
							"booker2": 2,
							"booker3": 1,
						},
					},
				},
			},
		},
	}
	memStore.StoreBatch(saveBucket,
		&define.FlowPack{
			Index: 1,
			Peers: map[string]uint32{
				"booker1": 2,
				"booker2": 1,
				"booker4":1,
			},
		},
	)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := memStore.UpdateBatch(tt.args.bucket, tt.args.batch...); (err != nil) != tt.wantErr {
				t.Errorf("UpdateBatch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInMemStore_DeleteRange(t *testing.T) {
	memStore.StoreBatch(saveBucket,
		&define.FlowPack{
			Index: 1,
			Peers: map[string]uint32{},
		},
		&define.FlowPack{
			Index: 2,
			Peers: map[string]uint32{},
		},
		&define.FlowPack{
			Index: 3,
			Peers: map[string]uint32{},
		},
		&define.FlowPack{
			Index: 100,
			Peers: map[string]uint32{},
		},
	)
	min := memStore.GetMin(saveBucket)
	max := memStore.GetMax(saveBucket)

	type args struct {
		bucket string
		min    iface.StoreItem
		max    iface.StoreItem
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:"TestInMemStore_DeleteRange",
			args:args{
				bucket: saveBucket,
				min:    min,
				max:    max,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := memStore.DeleteRange(tt.args.bucket, tt.args.min, tt.args.max); (err != nil) != tt.wantErr {
				t.Errorf("DeleteRange() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInMemStore_Delete(t *testing.T) {
	type args struct {
		bucket string
		items  []iface.StoreItem
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:"TestInMemStore_Delete",
			args:args{
				bucket: saveBucket,
				items:  []iface.StoreItem{&define.FlowPack{
					Index:      1,
				}},
			},
		},
	}

	memStore.StoreBatch(saveBucket,
		&define.FlowPack{
			Index: 1,
			Peers: map[string]uint32{},
		},
		&define.FlowPack{
			Index: 2,
			Peers: map[string]uint32{},
		},
		&define.FlowPack{
			Index: 3,
			Peers: map[string]uint32{},
		},
		&define.FlowPack{
			Index: 100,
			Peers: map[string]uint32{},
		},
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := memStore.Delete(tt.args.bucket, tt.args.items...); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}

			if got, err := memStore.GetBatch(tt.args.bucket, 100); err != nil {
				t.Errorf("GetBatch error = %v", err)
			} else {
				t.Logf("GotBatch num = %v", len(got))
			}
		})
	}
}