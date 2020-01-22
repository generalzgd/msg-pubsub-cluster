/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: subscribeflag_test.go.go
 * @time: 2020/1/3 1:00 下午
 * @project: msgsubscribesvr
 */

package define

import (
	`testing`

	`github.com/generalzgd/msg-subscriber/iproto`
)

var (
	p = NewSubscribeFlagMap()
)

func TestSubscribeFlagMap_Put(t *testing.T) {
	type args struct {
		nodeId string
		key    string
		ids    []int
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "TestSubscribeFlagMap_Put_1",
			args: args{
				nodeId: "node1",
				key:    "a",
				ids: []int{
					3846,
				},
			},
		},
		{
			name: "TestSubscribeFlagMap_Put_2",
			args: args{
				nodeId: "node2",
				key:    "b",
				ids: []int{
					3846,
				},
			},
		},
		{
			name: "TestSubscribeFlagMap_Put_3",
			args: args{
				nodeId: "node1",
				key:    "c",
				ids: []int{
					3846,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p.Put(tt.args.nodeId, tt.args.key, tt.args.ids)
		})
	}
}

func TestSubscribeFlagMap_Del(t *testing.T) {
	TestSubscribeFlagMap_Put(t)

	type args struct {
		nodeId string
		key    string
		ids    []int
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "TestSubscribeFlagMap_Del_1",
			args: args{
				nodeId: "node1",
				key:    "a",
				ids: []int{
					3846,
				},
			},
		},
		{
			name: "TestSubscribeFlagMap_Del_2",
			args: args{
				nodeId: "node2",
				key:    "a",
				ids: []int{
					3846,
				},
			},
		},
		{
			name: "TestSubscribeFlagMap_Del_3",
			args: args{
				nodeId: "node1",
				key:    "c",
				ids: []int{
					3846,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p.Del(tt.args.nodeId, tt.args.key, tt.args.ids)
		})
	}
}

func TestSubscribeFlagMap_Has(t *testing.T) {
	TestSubscribeFlagMap_Put(t)

	type args struct {
		id int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "TestSubscribeFlagMap_Has",
			args: args{
				id: 3846,
			},
			want: true,
		},
		{
			name: "TestSubscribeFlagMap_Has_2",
			args: args{
				id: 3300,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := p.Has(tt.args.id); got != tt.want {
				t.Errorf("SubscribeFlagMap.Has() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubscribeFlagMap_GetSubscribedCount(t *testing.T) {
	TestSubscribeFlagMap_Put(t)

	type args struct {
		id int
	}
	tests := []struct {
		name string
		args args
		want uint32
	}{
		// TODO: Add test cases.
		{
			name: "TestSubscribeFlagMap_GetSubscribedCount",
			args: args{
				id: 3846,
			},
			want: 3,
		},
		{
			name: "TestSubscribeFlagMap_GetSubscribedCount_2",
			args: args{
				id: 3300,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := p.GetSubscribedCount(tt.args.id); got != tt.want {
				t.Errorf("SubscribeFlagMap.GetSubscribedCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubscribeFlagMap_Set(t *testing.T) {
	TestSubscribeFlagMap_Put(t)

	type args struct {
		nodesMap map[string]*iproto.NodeMap
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "TestSubscribeFlagMap_Set",
			args: args{
				nodesMap: map[string]*iproto.NodeMap{
					"node1": {
						Data: map[string]*iproto.IdList{
							"a": {
								Ids: []uint32{3846},
							},
							"c": {
								Ids: []uint32{3846},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p.Set(tt.args.nodesMap)
		})
	}
}

func TestSubscribeFlagMap_Copy(t *testing.T) {
	TestSubscribeFlagMap_Set(t)

	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{
			name: "TestSubscribeFlagMap_Copy",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := p.Copy(); len(got) > 0 {
				t.Logf("SubscribeFlagMap.Copy() = %v", got)
			}
		})
	}
}
