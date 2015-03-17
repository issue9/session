// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package stores

import (
	"testing"
	"time"

	"github.com/issue9/assert"
	sess "github.com/issue9/session"
)

var _ sess.Store = &memory{}

// 声明两行测试数据。
var (
	testData1 = map[interface{}]interface{}{
		"10": 10,
		"11": 11,
	}

	testData2 = map[interface{}]interface{}{
		"20": 20,
		"21": 21,
	}
)

func TestMemory(t *testing.T) {
	a := assert.New(t)

	store := NewMemory(100)
	a.NotNil(store)

	// 添加一个数据
	a.NotError(store.Save("testData1", testData1))
	a.Equal(1, len(store.items))

	// Delete,删除一个不存在的数据，不应该发生错误
	a.NotError(store.Delete("non"))

	// Delete,删除添加的数据
	a.NotError(store.Delete("testData1"))
	a.Equal(0, len(store.items))

	// 添加两条数据
	a.NotError(store.Save("testData1", testData1))
	a.Equal(1, len(store.items))
	a.NotError(store.Save("testData2", testData2))
	a.Equal(2, len(store.items))

	// 测试正常状态的Get
	mapped, err := store.Get("testData1")
	a.NotError(err).NotNil(mapped)
	a.Equal(mapped, testData1)

	// 测试Get()一个不存在的数据。
	mapped, err = store.Get("non")
	a.NotError(err).Equal(0, len(mapped))

	// Free
	a.NotError(store.Save("testData1", testData1))
	a.NotError(store.Save("testData2", testData2))
	a.Equal(2, len(store.items))
	a.NotError(store.Close())
	a.Equal(0, len(store.items))
}

func TestMemory_StartGC(t *testing.T) {
	a := assert.New(t)

	// 2秒后开始执行GC
	store := NewMemory(2)
	a.NotNil(store)

	// 添加两条数据
	a.NotError(store.Save("testData1", testData1))
	a.NotError(store.Save("testData2", testData2))
	a.Equal(2, len(store.items))

	store.StartGC()
	a.Equal(2, len(store.items))
	time.Sleep(time.Second) // 延时1秒，数据还在
	a.Equal(2, len(store.items))
	time.Sleep(time.Second) // 再延时1秒，数据应该没了
	a.Equal(0, len(store.items))

	a.NotError(store.Close())
}
