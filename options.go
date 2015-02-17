// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package session

import (
	"net/http"
	"net/url"
	"time"
)

type Options struct {
	store    Store
	lifetime int
	cookie   *http.Cookie
}

func NewOptions(store Store, lifetime int, sessionIDName string, secure bool) *Options {
	return &Options{
		store:    store,
		lifetime: lifetime,
		cookie: &http.Cookie{
			Name:     sessionIDName,
			Secure:   secure,
			HttpOnly: true,
			Path:     "/",
		},
	}
}

// 设置相应的cookie值
func (opt *Options) setCookie(w http.ResponseWriter, sessid string, maxAge int) {
	opt.cookie.Value = url.QueryEscape(sessid)
	opt.cookie.MaxAge = maxAge
	// NOTE:ie8以下只支持Expires而不支持max_age。http1.0只有只有expires，
	// 而在http1.1中expires属于废弃的属性，max-age才是正规的。
	opt.cookie.Expires = time.Now().Add(time.Second * time.Duration(maxAge))

	http.SetCookie(w, opt.cookie)
}
