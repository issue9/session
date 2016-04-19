// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package session

import (
	"errors"
	"net/http"
	"sync"
)

// Session操作接口。
type Session struct {
	sync.Mutex

	manager *Manager
	id      string
	items   map[interface{}]interface{}
}

// 获取指定键名对应的值，found表示该值是否存在。
func (sess *Session) Get(key interface{}) (interface{}, bool) {
	sess.Lock()
	val, found := sess.items[key]
	sess.Unlock()
	return val, found
}

// 获取值，若键名对应的值不存在，则返回defVal。
func (sess *Session) MustGet(key, defVal interface{}) interface{} {
	sess.Lock()
	val, found := sess.items[key]
	sess.Unlock()

	if !found {
		return defVal
	}
	return val
}

// 添加或是设置值。
func (sess *Session) Set(key, val interface{}) {
	sess.Lock()
	sess.items[key] = val
	sess.Unlock()
}

// 指定的键值是否存在。
func (sess *Session) Exists(key interface{}) bool {
	sess.Lock()
	_, found := sess.items[key]
	sess.Unlock()

	return found
}

// 当前session的sessionid
func (sess *Session) ID() string {
	return sess.id
}

// 关闭当前的Session，相当于按顺序执行Session.Save()和Session.Free()。
func (sess *Session) Close(w http.ResponseWriter, r *http.Request) error {
	if err := sess.Save(w, r); err != nil {
		return err
	}

	return sess.Free(w, r)
}

// 释放当前的Session空间，但依然存在于Store中。
// 之后Session.Get等操作数据的函数将不在可用。
// 若需要同时从Store中去除，请执行Store.Delete()方法。
func (sess *Session) Free(w http.ResponseWriter, r *http.Request) error {
	sess.manager.provider.Delete(w, r)

	// 清空数据。
	sess.Lock()
	sess.items = nil
	sess.manager = nil
	sess.Unlock()

	return nil
}

// 保存当前的Session值到Store中。
// Session中的数据依然存在，可以继续使用Get()等函数获取数据。
func (sess *Session) Save(w http.ResponseWriter, r *http.Request) error {
	if sess.items == nil {
		return errors.New("数据已经被释放。")
	}

	return sess.manager.store.Save(sess.ID(), sess.items)
}
