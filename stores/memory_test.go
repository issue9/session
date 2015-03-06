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

	store := NewMemory()
	a.NotNil(store)

	// 添加一个数据
	a.NotError(store.Save("testData1", testData1))
	a.Equal(1, len(store.items))

	// 删除一个不存在的数据，不应该发生错误
	a.NotError(store.Delete("non"))

	// 删除添加的数据
	a.NotError(store.Delete("testData1"))
	a.Equal(0, len(store.items))

	// 添加两条数据
	a.NotError(store.Save("testData1", testData1))
	a.Equal(1, len(store.items))
	a.NotError(store.Save("testData2", testData2))
	a.Equal(2, len(store.items))

	// 测试Get
	mapped, err := store.Get("testData1")
	a.NotError(err).NotNil(mapped)
	a.Equal(mapped, testData1)

	// GC
	store.GC(2) // 生存时间为2秒
	a.Equal(2, len(store.items))
	// 休眠两秒后，所有数据都应该已经过期。
	time.Sleep(time.Second * 2)
	store.GC(2)
	a.Equal(0, len(store.items))

	// Free
	a.NotError(store.Save("testData1", testData1))
	a.NotError(store.Save("testData2", testData2))
	a.Equal(2, len(store.items))
	a.NotError(store.Free())
	a.Equal(0, len(store.items))
}
