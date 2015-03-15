// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package options

import (
	"net/http"
	"net/url"
	"time"

	"github.com/issue9/session"
)

// session操作的一些设置项。
// 目前sessionid保存于cookie中，cookie的设置都是通过Cookie完成的。
type cookie struct {
	store    session.Store
	lifetime int // 生存周期，单位为秒。
	ticker   *time.Ticker
	cookie   *http.Cookie
}

// 声明一个新的Cookie实例。
//
// store：该实例对应的Store接口，同时还会开始执行store.GC()；
// lifetime：Session的生存周期，单位为秒；
// sessionIDName sessionid在cookie中的名称；
// path,domain,secure也分别对应cookie中相应的值。
func NewCookie(store session.Store, lifetime int, sessionIDName, path, domain string, secure bool) session.Options {
	ticker := time.NewTicker(time.Duration(lifetime) * time.Second)
	go func() {
		for range ticker.C {
			store.GC(lifetime)
		}
	}()

	return &cookie{
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

// Options.Init()
func (c *cookie) Init(w http.ResponseWriter, req *http.Request) (sessID string, err error) {
	cookie, err := req.Cookie(c.cookie.Name)

	if err != nil || len(cookie.Value) == 0 { // 不存在，产生新的
		if sessID, err = sessionID(); err != nil {
			return sessID, err
		}
	} else { // 从Cookie中获取sessionid值。
		if sessID, err = url.QueryUnescape(cookie.Value); err != nil {
			return sessID, err
		}
	}

	c.cookie.Value = url.QueryEscape(sessID)
	c.cookie.MaxAge = c.lifetime
	// NOTE:ie8以下只支持Expires而不支持max_age；http1.0只有只有expires，
	// 而在http1.1中expires属于废弃的属性，max-age才是正规的。
	c.cookie.Expires = time.Now().Add(time.Second * time.Duration(c.lifetime))
	http.SetCookie(w, c.cookie)

	return sessID, nil
}

// Options.Delete()
func (c *cookie) Delete(w http.ResponseWriter, req *http.Request) error {
	c.cookie.MaxAge = -1
	http.SetCookie(w, c.cookie)

	return nil
}

// Options.Store()
func (c *cookie) Store() session.Store {
	return c.store
}

// 关闭Cookie及释放与之关联的Store，也会正常关闭Store.GC()的goroutinue。
func (c *cookie) Close() error {
	c.ticker.Stop()
	return c.store.Free()
}
