// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package session

import (
	"net/http"

	"github.com/issue9/session/types"
)

// session管理。
type Manager struct {
	store    types.Store
	provider types.Provider
}

// 声明一个Manager实例。
func New(store types.Store, prv types.Provider) *Manager {
	store.StartGC()

	return &Manager{
		store:    store,
		provider: prv,
	}
}

// 关闭，会自动释放关联的Store内容，即会删除所有的session数据。
func (mgr *Manager) Close() error {
	return mgr.store.Close()
}

// 获取与当前请求相关联的session数据。
// 在一个Session中，不能多次调用Start()。
// 当然也可以把获取的Session实例保存到Context等实例中，方便之后获取。
func (mgr *Manager) Start(w http.ResponseWriter, r *http.Request) (*Session, error) {
	sessID, err := mgr.provider.Get(w, r)
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
