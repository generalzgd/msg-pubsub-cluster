/**
 * @version: 1.0.0
 * @author: generalzgd
 * @license: LGPL v3
 * @contact: general_zgd@163.com
 * @site: github.com/generalzgd
 * @software: GoLand
 * @file: subscribemap_test.go.go
 * @time: 2019/12/24 10:20 上午
 * @project: packagesubscribesvr
 */

package define

import (
	`reflect`
	`sync`
	`testing`
	`time`

	`github.com/generalzgd/msg-subscriber/iface`
)

func TestSubscribeMap_mutex(t *testing.T) {

	p := NewSubscribeMap()
	outwg := sync.WaitGroup{}

	for loop := 0; loop < 5; loop++ {
		outwg.Add(2)
		go func() {
			defer outwg.Done()

			wg := sync.WaitGroup{}
			begin := time.Now()

			for i := 0; i < 1000; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					for i := 0; i < 10000; i++ {
						p.Put(nil, &SubscribeInfo{})
					}
				}()
			}
			wg.Wait()
			t.Logf("loop %d rwmutex write use time: %v", loop, time.Since(begin))
		}()

		//
		go func() {
			defer outwg.Done()

			wg := sync.WaitGroup{}
			begin := time.Now()

			for i := 0; i < 1000; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					for i := 0; i < 10000; i++ {
						p.GetConsumer(0)
					}
				}()
			}
			wg.Wait()
			t.Logf("loop %d rwmutex read use time: %v", loop, time.Since(begin))
		}()

		outwg.Wait()
	}
	/*
			statistic:
		subscribemap_test.go:46: loop 0 mutex write use time: 7.263519003s
		subscribemap_test.go:68: loop 0 mutex read use time: 7.263529744s
		subscribemap_test.go:68: loop 1 mutex read use time: 8.090125773s
		subscribemap_test.go:46: loop 1 mutex write use time: 8.090135579s
		subscribemap_test.go:68: loop 2 mutex read use time: 6.073017783s
		subscribemap_test.go:46: loop 2 mutex write use time: 6.072956168s
		subscribemap_test.go:68: loop 3 mutex read use time: 4.142696456s
		subscribemap_test.go:46: loop 3 mutex write use time: 4.142811145s
		subscribemap_test.go:46: loop 4 mutex write use time: 4.492920379s
		subscribemap_test.go:68: loop 4 mutex read use time: 4.493133377s

		subscribemap_test.go:68: loop 0 rwmutex read use time: 710.911503ms
		subscribemap_test.go:46: loop 0 rwmutex write use time: 2.022427732s
		subscribemap_test.go:68: loop 1 rwmutex read use time: 713.654163ms
		subscribemap_test.go:46: loop 1 rwmutex write use time: 1.860055015s
		subscribemap_test.go:68: loop 2 rwmutex read use time: 714.820677ms
		subscribemap_test.go:46: loop 2 rwmutex write use time: 1.914173567s
		subscribemap_test.go:68: loop 3 rwmutex read use time: 727.135025ms
		subscribemap_test.go:46: loop 3 rwmutex write use time: 1.838603784s
		subscribemap_test.go:68: loop 4 rwmutex read use time: 723.596415ms
		subscribemap_test.go:46: loop 4 rwmutex write use time: 1.854431875s

			compare:
		   mutex read: 7.2  8.0  6.0  4.1  4.4
		 rwmutex read: 0.71 0.71 0.71 0.72 0.72
		  mutex write: 7.2  8.0  6.0  4.1  4.4
		rwmutex write: 2.0  1.8  1.9  1.8  1.8

			result:
			rwmutex read is 10 times faster than mutex read
			rwmutex write is 4 times faster than mutex write
	*/
}

var (
	mp = NewSubscribeMap()
)

func TestSubscribeMap_Put(t *testing.T) {
	type args struct {
		ids []int
		v   iface.IConsumer
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		// TODO: Add test cases.
		{
			name: "TestSubscribeMap_Put_1", // 加入
			args: args{
				ids: []int{1, 2, 3},
				v: &SubscribeInfo{
					ConsumerKey: "a",
				},
			},
			want: []int{1, 2, 3},
		},
		{
			name: "TestSubscribeMap_Put_2", // 已经存在，加入失败
			args: args{
				ids: []int{1},
				v: &SubscribeInfo{
					ConsumerKey: "a",
				},
			},
			want: []int{1},
		},
		{
			name: "TestSubscribeMap_Put_3", // 部分加入成功
			args: args{
				ids: []int{1, 4},
				v: &SubscribeInfo{
					ConsumerKey: "a",
				},
			},
			want: []int{1, 4},
		},
		{
			name: "TestSubscribeMap_Put_4", // 同杨的协议，不同的订阅者
			args: args{
				ids: []int{1, 4},
				v: &SubscribeInfo{
					ConsumerKey: "b",
				},
			},
			want: []int{1, 4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mp.Put(tt.args.ids, tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Put() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubscribeMap_Del(t *testing.T) {
	type args struct {
		ids []int
		v   iface.IConsumer
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		// TODO: Add test cases.
		{
			name: "TestSubscribeMap_Del_1", // 正常删除
			args: args{
				ids: []int{1},
				v: &SubscribeInfo{
					ConsumerKey: "a",
				},
			},
			want: []int{1},
		},
		{
			name: "TestSubscribeMap_Del_1", // 删除未注册过的ID
			args: args{
				ids: []int{5},
				v: &SubscribeInfo{
					ConsumerKey: "a",
				},
			},
			want: []int{},
		},
		{
			name: "TestSubscribeMap_Del_1", // 删除未注册过的ID和正常的ID
			args: args{
				ids: []int{3, 4, 5},
				v: &SubscribeInfo{
					ConsumerKey: "b",
				},
			},
			want: []int{3, 4},
		},
	}
	mp.Put([]int{1, 2, 3, 4}, &SubscribeInfo{
		ConsumerKey: "a",
	})
	mp.Put([]int{3, 4}, &SubscribeInfo{
		ConsumerKey: "b",
	})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mp.Del(tt.args.ids, tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Del() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubscribeMap_GetConsumer(t *testing.T) {
	type args struct {
		id int
	}
	tests := []struct {
		name  string
		args  args
		want  []iface.IConsumer
		want1 bool
	}{
		// TODO: Add test cases.
		{
			name: "TestSubscribeMap_GetConsumer_1", // 正常获取已注册的消费者
			args: args{id: 1},
			want: []iface.IConsumer{
				&SubscribeInfo{
					ConsumerKey: "a",
				},
				&SubscribeInfo{
					ConsumerKey: "b",
				},
			},
			want1: true,
		},
		{
			name:  "TestSubscribeMap_GetConsumer_2", // 获取未注册的消费者
			args:  args{id: 5},
			want:  nil,
			want1: false,
		},
	}

	mp.Put([]int{1, 2, 3, 4}, &SubscribeInfo{ConsumerKey: "a"})
	mp.Put([]int{1, 2}, &SubscribeInfo{ConsumerKey: "b"})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := mp.GetConsumer(tt.args.id)
			t.Logf("GetConsumer() got = %v, want %v", got, tt.want)
			if got1 != tt.want1 {
				t.Errorf("GetConsumer() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestSubscribeMap_GetConsumerByKey(t *testing.T) {
	type args struct {
		keys []string
	}
	tests := []struct {
		name string
		args args
		want []iface.IConsumer
	}{
		// TODO: Add test cases.
		{
			name: "TestSubscribeMap_GetConsumerByKey_1", // 正常获取
			args: args{keys: []string{"a"}},
			want: []iface.IConsumer{
				&SubscribeInfo{ConsumerKey: "a"},
			},
		},
		{
			name: "TestSubscribeMap_GetConsumerByKey_2", // 获取不存在的
			args: args{keys: []string{"c"}},
			want: []iface.IConsumer{},
		},
		{
			name: "TestSubscribeMap_GetConsumerByKey_3", // 获取多个，包含不存在的key
			args: args{keys: []string{"a", "b", "c"}},
			want: []iface.IConsumer{
				&SubscribeInfo{ConsumerKey: "a"},
				&SubscribeInfo{ConsumerKey: "b"},
			},
		},
	}

	mp.Put([]int{1, 2}, &SubscribeInfo{ConsumerKey: "a"})
	mp.Put([]int{3, 4}, &SubscribeInfo{ConsumerKey: "b"})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mp.GetConsumerByKey(tt.args.keys...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetConsumerByKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubscribeMap_HasBooked(t *testing.T) {
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
			name: "TestSubscribeMap_HasBooked_1",
			args: args{id: 1},
			want: true,
		},
		{
			name: "TestSubscribeMap_HasBooked_2",
			args: args{id: 103},
			want: false,
		},
	}

	mp.Put([]int{1, 2}, &SubscribeInfo{ConsumerKey: "a"})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mp.HasBooked(tt.args.id); got != tt.want {
				t.Errorf("HasBooked() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubscribeMap_GetIdsByConsumer(t *testing.T) {
	type args struct {
		tar iface.IConsumer
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		// TODO: Add test cases.
		{
			name: "TestSubscribeMap_GetIdsByConsumer_1",
			args: args{tar: &SubscribeInfo{
				ConsumerKey: "aa",
			}},
			want: []int{1, 2},
		},
		{
			name: "TestSubscribeMap_GetIdsByConsumer_2",
			args: args{tar: &SubscribeInfo{
				ConsumerKey: "bb",
			}},
			want: nil,
		},
	}

	mp.Put([]int{1, 2}, &SubscribeInfo{ConsumerKey: "aa"})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mp.GetIdsByConsumer(tt.args.tar)
			t.Logf("GetIdsByConsumer() got = %v", got)
		})
	}
}
