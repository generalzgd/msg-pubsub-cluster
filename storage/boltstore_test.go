/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: boltstore_test.go.go
 * @time: 2019/12/24 5:34 下午
 * @project: packagesubscribesvr
 */

package storage

import (
	`encoding/binary`
	`os`
	`path/filepath`
	`testing`

	`github.com/generalzgd/msg-subscriber/codec`
	`github.com/generalzgd/msg-subscriber/codec/body`
	`github.com/generalzgd/msg-subscriber/codec/head`
	`github.com/generalzgd/msg-subscriber/define`
	`github.com/generalzgd/msg-subscriber/iface`
)

var (
	boltBucket      = "bolttest"
	db              *BoltStore
	flowPackFactory func(string) iface.StoreItem
	headDecoder     head.Decoder
	bodyDecoder     body.Decoder
)

func init() {
	headDecoder = head.NewDecoder(22, 2, 0, 4, binary.LittleEndian)
	bodyDecoder = body.NewDecoder(24, binary.LittleEndian)
	flowPackFactory = func(str string) iface.StoreItem {
		return define.NewFlowPack(headDecoder, bodyDecoder)
	}
	path := filepath.Join(filepath.Dir(os.Args[0]), "bolt.db")
	db = NewBoltStore(path, flowPackFactory, boltBucket)
}

func TestBoltStore_GetUint64(t *testing.T) {
	bucket, key := "test", "index"
	path := filepath.Join(filepath.Dir(os.Args[0]), "tt.bolt")
	db := NewBoltStore(path, func(str string) iface.StoreItem {
		return &define.IndexPack{}
	}, bucket)

	if _, err := db.SetUnit64(bucket, key, 1234); err != nil {
		t.Log("set error.", err)
		return
	}
	if got, err := db.GetUint64(bucket, key); err != nil {
		t.Log("set error.", err)
		return
	} else {
		t.Log("got result:", got)
	}

}

func TestBoltStore_Store(t *testing.T) {
	pk := define.NewFlowPack(headDecoder, bodyDecoder)
	pk.Index = 11
	pk.Pack = codec.NewDataPack(nil, headDecoder, bodyDecoder)

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
			name: "TestBoltStore_Store",
			args: args{
				bucket: boltBucket,
				val:    pk,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := db.Store(tt.args.bucket, tt.args.val); (err != nil) != tt.wantErr {
				t.Errorf("Store() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBoltStore_StoreBatch(t *testing.T) {
	pk := define.NewFlowPack(headDecoder, bodyDecoder)
	pk.Index = 12
	pk.Pack = codec.NewDataPack(nil, headDecoder, bodyDecoder)

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
			name: "TestBoltStore_StoreBatch",
			args: args{
				bucket: boltBucket,
				batch: []iface.StoreItem{
					pk,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := db.StoreBatch(tt.args.bucket, tt.args.batch...); (err != nil) != tt.wantErr {
				t.Errorf("StoreBatch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBoltStore_GetBatch(t *testing.T) {
	type args struct {
		bucket string
		limit  int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "TestBoltStore_GetBatch",
			args: args{
				bucket: boltBucket,
				limit:  100,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := db.GetBatch(tt.args.bucket, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBatch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("GetBatch() num = %v", len(got))
		})
	}
}

func TestBoltStore_Delete(t *testing.T) {
	pk1 := define.NewFlowPack(headDecoder, bodyDecoder)
	pk1.Index = 11
	pk1.Pack = codec.NewDataPack(nil, headDecoder, bodyDecoder)

	pk2 := define.NewFlowPack(headDecoder, bodyDecoder)
	pk2.Index = 12
	pk2.Pack = codec.NewDataPack(nil, headDecoder, bodyDecoder)

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
			name: "TestBoltStore_Delete",
			args: args{
				bucket: boltBucket,
				items: []iface.StoreItem{
					pk1,
					pk2,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := db.Delete(tt.args.bucket, tt.args.items...); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBoltStore_UpdateBatch(t *testing.T) {
	pk1 := define.NewFlowPack(headDecoder, bodyDecoder)
	pk1.Index = 12
	pk1.Peers = map[string]uint32{
		"booker1": 2,
		"booker2": 1,
		"booker3": 1,
	}
	pk1.Pack = codec.NewDataPack(nil, headDecoder, bodyDecoder)

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
			name: "TestBoltStore_UpdateBatch_1",
			args: args{
				bucket: boltBucket,
				batch: []iface.StoreItem{
					pk1,
				},
			},
		},
		{
			name: "TestBoltStore_UpdateBatch_2",
			args: args{
				bucket: boltBucket,
				batch: []iface.StoreItem{
					pk1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := db.UpdateBatch(tt.args.bucket, tt.args.batch...); (err != nil) != tt.wantErr {
				t.Errorf("UpdateBatch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBoltStore_DeleteRange(t *testing.T) {
	min := db.GetMin(boltBucket)
	max := db.GetMax(boltBucket)
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
			name: "TestBoltStore_DeleteRange",
			args: args{
				bucket: boltBucket,
				min:    min,
				max:    max,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := db.DeleteRange(tt.args.bucket, tt.args.min, tt.args.max); (err != nil) != tt.wantErr {
				t.Errorf("DeleteRange() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
