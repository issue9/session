// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package session

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/issue9/assert"
)

// 测试Session的存储功能
func TestSessionAccess(t *testing.T) {
	a := assert.New(t)

	// 声明Options实例。
	store := newTestStore()
	opt := NewOptions(store, 10, "gosession", "/", "localhost", true)
	a.NotNil(opt)
	defer func() {
		a.NotError(opt.Close())
	}()

	h := func(w http.ResponseWriter, req *http.Request) {
		sess, err := Start(opt, w, req)
		a.NotError(err).NotNil(sess)

		// 通过多次调用Start()，返回的数据应该是不相同的。
		sess1, err := Start(opt, w, req)
		a.NotError(err).NotNil(sess1)
		a.NotEqual(sess1, sess)

		// 不存在的值
		val, found := sess.Get("nil")
		a.False(found).Nil(val)

		val = sess.MustGet("nil", "default")
		a.False(sess.Exists("nil")).Equal(val, "default")

		// MustGet()不应该记住值，所以此处还是空值。
		val, found = sess.Get("nil")
		a.False(found).Nil(val)

		// 设置值，并可以正确取回。
		sess.Set(5, "5")
		val, found = sess.Get(5)
		a.True(sess.Exists(5)).True(found).Equal(val, "5")
		val = sess.MustGet(5, "10")
		a.Equal(val, "5")

		// 键值可以为nil
		sess.Set(nil, "nil")
		val, found = sess.Get(nil)
		a.True(sess.Exists(nil)).True(found).Equal(val, "nil")
		sess.Set(nil, 5)
		val, found = sess.Get(nil)
		a.True(found).Equal(val, 5)

		// 添加了2个值。nil和5
		a.Equal(2, len(sess.items))

		// 此时store.items的长度应该为0
		a.Equal(0, len(store.items))

		// 保存数据到store
		a.NotError(sess.Close(w, req))
		a.Equal(1, len(store.items)) // store.items的长度变更为1
		// 此时，应该能通过sess.ID()正确找到该元素。
		item, found := store.items[sess.ID()]
		a.True(found).NotNil(item)
		a.Equal(0, len(sess.items)) // Close()，sess.items数据将被清空。
	}
	srv := httptest.NewServer(http.HandlerFunc(h))
	a.NotNil(srv)
	defer srv.Close()

	response, err := http.Get(srv.URL)
	a.NotError(err).NotNil(response)
}
