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
)

func TestSubscribeFlagMap_Put(t *testing.T) {
	type args struct {
		cmdIds []int
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "TestSubscribeFlagMap_Put",
			args: args{cmdIds: []int{1, 2, 3, 4, 5, 1, 2}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewSubscribeFlagMap()
			p.Put(tt.args.cmdIds...)
		})
	}
}

func TestSubscribeFlagMap_Del(t *testing.T) {
	type args struct {
		cmdIds []int
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "TestSubscribeFlagMap_Del",
			args: args{cmdIds: []int{1, 2}},
		},
	}
	p := NewSubscribeFlagMap()
	p.Put(1, 2, 3, 4, 5, 6)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p.Del(tt.args.cmdIds...)
		})
	}
}

func TestSubscribeFlagMap_Has(t *testing.T) {
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
			args: args{id: 2},
			want: true,
		},
	}
	p := NewSubscribeFlagMap()
	p.Put(1, 2, 3, 4, 5, 6)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := p.Has(tt.args.id); got != tt.want {
				t.Errorf("Has() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubscribeFlagMap_GetSubscribedCount(t *testing.T) {
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
			name: "TestSubscribeFlagMap_GetSubscribedCount_1",
			args: args{id: 1},
			want: 1,
		},
		{
			name: "TestSubscribeFlagMap_GetSubscribedCount_2",
			args: args{id: 7},
			want: 0,
		},
	}

	p := NewSubscribeFlagMap()
	p.Put(1, 2, 3, 4, 5, 6)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := p.GetSubscribedCount(tt.args.id); got != tt.want {
				t.Errorf("GetSubscribedCount() = %v, want %v", got, tt.want)
			}
		})
	}
}
