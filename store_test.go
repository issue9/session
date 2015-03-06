// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package session

import (
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
