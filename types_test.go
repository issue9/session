// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package session

import (
	"net/http"
	"net/url"
	"strconv"
	"time"
)

var _ Store = &testStore{}

// 测试用的Store接口实现。
type testStore struct {
	items map[string]*testSession
}

type testSession struct {
	accessed time.Time
	items    map[interface{}]interface{}
}

func newTestStore() *testStore {
	return &testStore{
		items: map[string]*testSession{},
	}
}

// Store.Delete()
func (store *testStore) Delete(sessID string) error {
	delete(store.items, sessID)
	return nil
}

// Store.Get()
func (store *testStore) Get(sessID string) (map[interface{}]interface{}, error) {
	if item, found := store.items[sessID]; found {
		return item.items, nil
	}

	return make(map[interface{}]interface{}, 0), nil
}

// Store.Save()
func (store *testStore) Save(sessID string, items map[interface{}]interface{}) error {
	store.items[sessID] = &testSession{
		accessed: time.Now(),
		items:    items,
	}
	return nil
}

// Store.GC()
func (store *testStore) GC(maxAge int) error {
	d := time.Now().Add(time.Duration(maxAge))

	for k, v := range store.items {
		if v.accessed.Before(d) { // 过期，则删除
			delete(store.items, k)
		}
	}
	return nil
}

// Store.Free()
func (store *testStore) Free() error {
	store.items = nil
	return nil
}

var _ Options = &testOptions{}

// 测试用Options接口的实现。
type testOptions struct {
	store    Store
	lifetime int // 生存周期，单位为秒。
	ticker   *time.Ticker
	cookie   *http.Cookie
	count    int // 用于产生唯一ID
}

func newTestOptions(store Store, lifetime int, sessionIDName string) Options {
	ticker := time.NewTicker(time.Duration(lifetime) * time.Second)
	go func() {
		for range ticker.C {
			store.GC(lifetime)
		}
	}()

	return &testOptions{
		store:    store,
		lifetime: lifetime,
		ticker:   ticker,
		cookie: &http.Cookie{
			Name:     sessionIDName,
			Secure:   true,
			HttpOnly: true,
			Path:     "/",
			Domain:   "localhost",
		},
	}
}

// Options.Init()
func (opt *testOptions) Init(w http.ResponseWriter, req *http.Request) (sessID string, err error) {
	cookie, err := req.Cookie(opt.cookie.Name)

	if err != nil || len(cookie.Value) == 0 { // 不存在，产生新的
		opt.count++
		sessID = "gosessionid:" + strconv.Itoa(opt.count)
	} else { // 从Cookie中获取sessionid值。
		if sessID, err = url.QueryUnescape(cookie.Value); err != nil {
			return sessID, err
		}
	}

	opt.cookie.Value = url.QueryEscape(sessID)
	opt.cookie.MaxAge = opt.lifetime
	http.SetCookie(w, opt.cookie)

	return sessID, nil
}

// Options.Delete()
func (opt *testOptions) Delete(w http.ResponseWriter, req *http.Request) error {
	opt.cookie.MaxAge = -1
	http.SetCookie(w, opt.cookie)

	return nil
}

// Options.Store()
func (opt *testOptions) Store() Store {
	return opt.store
}

// 关闭Cookie及释放与之关联的Store，也会正常关闭Store.GC()的goroutinue。
func (opt *testOptions) Close() error {
	opt.ticker.Stop()
	return opt.store.Free()
}
