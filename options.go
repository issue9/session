// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package session

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"
)

// session操作的一些设置项。
// 目前sessionid保存于cookie中，cookie的设置都是通过Options完成的。
type Options struct {
	store    Store
	lifetime int // 生存周期，单位为秒。
	ticker   *time.Ticker
	cookie   *http.Cookie
}

// 声明一个新的Options实例。
// store：该实例对应的Store接口；lifetime：Session的生存周期，单位为秒；
func NewOptions(store Store, lifetime int, sessionIDName, path, domain string, secure bool) *Options {
	ticker := time.NewTicker(time.Duration(lifetime) * time.Second)
	go func() {
		for range ticker.C {
			store.GC(lifetime)
		}
	}()

	return &Options{
		store:    store,
		lifetime: lifetime,
		ticker:   ticker,
		cookie: &http.Cookie{
			Name:     sessionIDName,
			Secure:   secure,
			HttpOnly: true,
			Path:     path,
			Domain:   domain,
		},
	}
}

// 设置cookie的sessid和maxage值。
func (opt *Options) setCookie(w http.ResponseWriter, sessid string, maxAge int) {
	if maxAge > 0 { // 若是删除属性，则不用计算该值。
		opt.cookie.Value = url.QueryEscape(sessid)
	}

	opt.cookie.MaxAge = maxAge
	// NOTE:ie8以下只支持Expires而不支持max_age；http1.0只有只有expires，
	// 而在http1.1中expires属于废弃的属性，max-age才是正规的。
	opt.cookie.Expires = time.Now().Add(time.Second * time.Duration(maxAge))

	http.SetCookie(w, opt.cookie)
}

// 根据当前的req获取相应的sessionid
func (opt *Options) getSessionID(req *http.Request) (string, error) {
	cookie, err := req.Cookie(opt.cookie.Name)
	if err == nil && cookie.Value != "" {
		return url.QueryUnescape(cookie.Value)
	}

	// 不存在，新建一个sessionid
	ret := make([]byte, 64)
	n, err := io.ReadFull(rand.Reader, ret)
	if n == 0 {
		return "", errors.New("未读取到随机数")
	}

	h := md5.New()
	h.Write(ret)
	return hex.EncodeToString(h.Sum(nil)), err
}

// 关闭Options及释放与之关联的Store
func (opt *Options) Close() error {
	opt.ticker.Stop()
	return opt.store.Free()
}
