// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package stores

import (
	"testing"
	"time"

	"github.com/issue9/assert"
)

func TestFile(t *testing.T) {
	a := assert.New(t)

	store, err := NewFile("./testdata")
	a.NotError(err).NotNil(store)

	// 添加一个数据
	a.NotError(store.Save("testData1", testData1))
	a.FileExists(store.dir + "testData1")

	// Delete,删除一个不存在的数据，不应该发生错误
	a.NotError(store.Delete("non"))

	// Delete,删除添加的数据
	a.NotError(store.Delete("testData1"))
	a.FileNotExists(store.dir + "testData1")

	// 添加两条数据
	a.NotError(store.Save("testData1", testData1))
	a.FileExists(store.dir + "testData1")
	a.NotError(store.Save("testData2", testData2))
	a.FileExists(store.dir + "testData2")

	// 测试正常状态的Get
	mapped, err := store.Get("testData1")
	a.NotError(err).NotNil(mapped)
	a.Equal(mapped, testData1)

	// 测试Get()一个不存在的数据。
	mapped, err = store.Get("non")
	a.NotError(err).Equal(0, len(mapped))

	// GC
	store.GC(2) // 生存时间为2秒
	a.FileExists(store.dir + "testData1")
	a.FileExists(store.dir + "testData2")
	// 休眠两秒后，所有数据都应该已经过期。
	time.Sleep(time.Second * 2)
	store.GC(2)
	a.FileNotExists(store.dir + "testData1")
	a.FileNotExists(store.dir + "testData2")

	// Free
	a.NotError(store.Save("testData1", testData1))
	a.NotError(store.Save("testData2", testData2))
	a.FileExists(store.dir + "testData1")
	a.FileExists(store.dir + "testData2")
	a.NotError(store.Free())
	a.FileNotExists(store.dir + "testData1")
	a.FileNotExists(store.dir + "testData2")
}
