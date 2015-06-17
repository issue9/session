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

func TestManager_Start(t *testing.T) {
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

		// 通过多次调用Start()，返回的数据应该是不相同的。
		sess1, err := mgr.Start(w, req)
		a.NotError(err).NotNil(sess1)
		a.NotEqual(sess1, sess)
	}
	srv := httptest.NewServer(http.HandlerFunc(h))
	a.NotNil(srv)
	defer srv.Close()

	response, err := http.Get(srv.URL)
	a.NotError(err).NotNil(response)
}
