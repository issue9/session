// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// session的内存存储模式
package stores

import (
	"sync"
	"time"
)

type memSession struct {
	accessed time.Time
	items    map[interface{}]interface{}
}

type memory struct {
	sync.Mutex

	items map[string]*memSession
}

// 返回一个实现session.Store接口的内存存储器。
func NewMemory() *memory {
	return &memory{
		items: map[string]*memSession{},
	}
}

// session.Store.Delete()
func (mem *memory) Delete(sessID string) error {
	mem.Lock()
	defer mem.Unlock()

	delete(mem.items, sessID)
	return nil
}

// session.Store.Get()
func (mem *memory) Get(sessID string) (map[interface{}]interface{}, error) {
	mem.Lock()
	defer mem.Unlock()

	if item, found := mem.items[sessID]; found {
		return item.items, nil
	}

	return make(map[interface{}]interface{}, 0), nil
}

// session.Store.Save()
func (mem *memory) Save(sessID string, items map[interface{}]interface{}) error {
	mem.Lock()
	defer mem.Unlock()

	mem.items[sessID] = &memSession{
		accessed: time.Now(),
		items:    items,
	}
	return nil
}

// session.Store.GC()
func (mem *memory) GC(maxAge int) error {
	d := time.Now().Add(-time.Second * time.Duration(maxAge))

	for k, v := range mem.items {
		if v.accessed.Before(d) { // v.accessed < (time.Now() - maxAge)
			delete(mem.items, k)
		}
	}
	return nil
}

// session.Store.Free()
func (mem *memory) Free() error {
	mem.items = nil
	return nil
}
