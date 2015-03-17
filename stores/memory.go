// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

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

	items    map[string]*memSession
	ticker   *time.Ticker
	lifetime time.Duration
}

// 返回一个实现session.Store接口的内存存储器。
//
// 内存存储器是不稳定的，随着程序中止或是实例被销毁，
// 相关的session数据也会随之销毁。
func NewMemory(lifetime int) *memory {
	return &memory{
		lifetime: time.Second * time.Duration(lifetime),
		items:    map[string]*memSession{},
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

// session.Store.StartGC()
func (mem *memory) StartGC() {
	gc := func() {
		d := time.Now().Add(-mem.lifetime)

		for k, v := range mem.items {
			if v.accessed.Before(d) { // v.accessed < (time.Now() - maxAge)
				delete(mem.items, k)
			}
		}
	}

	mem.ticker = time.NewTicker(mem.lifetime)
	go func() {
		for range mem.ticker.C {
			gc()
		}
	}()
}

// session.Store.Close()
func (mem *memory) Close() error {
	if mem.ticker != nil {
		mem.ticker.Stop()
	}

	mem.items = nil
	return nil
}
