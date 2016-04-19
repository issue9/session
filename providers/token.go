// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package providers

import (
	"net/http"
)

type token struct {
	name     string // token在报头中的名称
	lifetime int
}

func NewToken(name string, lifetime int) *token {
	return &token{
		name:     name,
		lifetime: lifetime,
	}
}

// session.Provider.Get()
func (t *token) Get(w http.ResponseWriter, req *http.Request) (sessID string, err error) {
	sessID = req.Header.Get(t.name)
	if len(sessID) == 0 {
		if sessID, err = sessionID(); err != nil {
			return "", err
		}
	}

	return sessID, nil
}

// session.Provider.Delete()
func (t *token) Delete(w http.ResponseWriter, req *http.Request) error {
	// token 由用户在客户端维持。服务端并不用做特殊处理。
	return nil
}
