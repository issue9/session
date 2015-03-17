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
	items    map[string]*testSession
	lifetime int
}

type testSession struct {
	accessed time.Time
	items    map[interface{}]interface{}
}

func newTestStore() *testStore {
	return &testStore{
		items:    map[string]*testSession{},
		lifetime: 10,
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

// Store.StartGC()
func (store *testStore) StartGC() {
	d := time.Now().Add(time.Duration(store.lifetime))

	for k, v := range store.items {
		if v.accessed.Before(d) { // 过期，则删除
			delete(store.items, k)
		}
	}
}

// Store.Close()
func (store *testStore) Close() error {
	store.items = nil
	return nil
}

var _ Provider = &testProvider{}

// 测试用Manager接口的实现。
type testProvider struct {
	lifetime int // 生存周期，单位为秒。
	cookie   *http.Cookie
	count    int // 用于产生唯一ID
}

func newTestManager(lifetime int, sessionIDName string) Provider {
	return &testProvider{
		lifetime: lifetime,
		cookie: &http.Cookie{
			Name:     sessionIDName,
			Secure:   true,
			HttpOnly: true,
			Path:     "/",
			Domain:   "localhost",
		},
	}
}

// Manager.Get()
func (mgr *testProvider) Get(w http.ResponseWriter, req *http.Request) (sessID string, err error) {
	cookie, err := req.Cookie(mgr.cookie.Name)

	if err != nil || len(cookie.Value) == 0 { // 不存在，产生新的
		mgr.count++
		sessID = "gosessionid:" + strconv.Itoa(mgr.count)
	} else { // 从Cookie中获取sessionid值。
		if sessID, err = url.QueryUnescape(cookie.Value); err != nil {
			return sessID, err
		}
	}

	mgr.cookie.Value = url.QueryEscape(sessID)
	mgr.cookie.MaxAge = mgr.lifetime
	http.SetCookie(w, mgr.cookie)

	return sessID, nil
}

// Manager.Delete()
func (mgr *testProvider) Delete(w http.ResponseWriter, req *http.Request) error {
	mgr.cookie.MaxAge = -1
	http.SetCookie(w, mgr.cookie)

	return nil
}
