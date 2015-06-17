// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package session

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/session/providers"
	"github.com/issue9/session/stores"
)

// 测试Session的存储功能
func TestSessionAccess1(t *testing.T) {
	a := assert.New(t)

	// 声明Manager实例。
	store := stores.NewMemory(10)
	prv := providers.NewCookie(10, "gosession", "/", "localhost", true)
	mgr := New(store, prv)
	a.NotNil(mgr)
	defer func() {
		a.NotError(mgr.Close())
	}()

	h := func(w http.ResponseWriter, req *http.Request) {
		sess, err := mgr.Start(w, req)
		a.NotError(err).NotNil(sess)

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

		// 保存数据到store
		a.NotError(sess.Close(w, req))

		// 此时，应该能通过sess.ID()正确找到该元素。
		item, err := store.Get(sess.ID())
		a.NotError(found).NotNil(item)
		a.Equal(0, len(sess.items)) // Close()，sess.items数据将被清空。
	}
	srv := httptest.NewServer(http.HandlerFunc(h))
	a.NotNil(srv)
	defer srv.Close()

	response, err := http.Get(srv.URL)
	a.NotError(err).NotNil(response)
}

// 测试多个store同时使用
func TestSessionAccess2(t *testing.T) {
	a := assert.New(t)

	s1 := stores.NewMemory(10)
	s2 := stores.NewMemory(10)
	p1 := providers.NewCookie(10, "gosession1", "/", "locahost", true)
	p2 := providers.NewCookie(10, "gosession2", "/", "locahost", true)
	mgr1 := New(s1, p1)
	mgr2 := New(s2, p2)
	defer mgr1.Close()
	defer mgr2.Close()

	h := func(w http.ResponseWriter, req *http.Request) {
		sess1, err := mgr1.Start(w, req)
		a.NotError(err).NotNil(sess1)

		sess2, err := mgr2.Start(w, req)
		a.NotError(err).NotNil(sess2)
		a.NotEqual(sess1, sess2)

		// 不存在的值
		val, found := sess1.Get("nil")
		a.False(found).Nil(val)
		val, found = sess2.Get("nil")
		a.False(found).Nil(val)

		// 仅设置了sess1，sess2应该不存在该值
		sess1.Set("1", 1)
		a.True(sess1.Exists("1")).False(sess2.Exists("1"))

		sess2.Set("2", 2)
		a.False(sess1.Exists("2")).True(sess2.Exists("2"))

		// 销毁Sess1，应该不影响sess2
		sess1.Close(w, req)
		a.False(sess1.Exists("1"))
		a.True(sess2.Exists("2"))
	}
	srv := httptest.NewServer(http.HandlerFunc(h))
	a.NotNil(srv)
	defer srv.Close()

	response, err := http.Get(srv.URL)
	a.NotError(err).NotNil(response)
}
