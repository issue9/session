// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package session

import (
	"net/http"
)

type Manager struct {
	store    Store
	provider Provider
}

func New(store Store, prv Provider) *Manager {
	store.StartGC()

	return &Manager{
		store:    store,
		provider: prv,
	}
}

// 关闭，会自动释放关联的Store内容
func (mgr *Manager) Close() error {
	return mgr.store.Close()
}

// 当然也可以把获取的Session实例保存到Context等实例中，方便之后获取。
// 在一个Session中，不能多次调用Start()。
func (mgr *Manager) Start(w http.ResponseWriter, req *http.Request) (*Session, error) {
	sessID, err := mgr.provider.Get(w, req)
	if err != nil {
		return nil, err
	}

	items, err := mgr.store.Get(sessID)
	if err != nil {
		return nil, err
	}

	return &Session{
		manager: mgr,
		id:      sessID,
		items:   items,
	}, nil
}
