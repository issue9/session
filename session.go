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
	"sync"
)

// Session操作接口
type Session struct {
	sync.Mutex

	id      string
	items   map[interface{}]interface{}
	options *Options
}

// 返回当前request的Session实例。
func NewSession(opt *Options, w http.ResponseWriter, req *http.Request) (*Session, error) {
	sess := &Session{
		id:      "",
		items:   make(map[interface{}]interface{}, 1),
		options: opt,
	}

	// 获取sessionid的值
	cookie, err := req.Cookie(sess.options.cookie.Name)
	if err != nil || cookie.Value == "" { // 不存在，新建一个sessionid
		sess.id, err = sessionID()
	} else {
		sess.id, err = url.QueryUnescape(cookie.Value)
	}

	if err != nil {
		return nil, err
	}

	opt.setCookie(w, sess.id, opt.lifetime)
	return sess, nil
}

// 获取指定键名对应的值，found表示该值是否存在。
func (sess *Session) Get(key interface{}) (val interface{}, found bool) {
	sess.Lock()
	defer sess.Unlock()

	val, found = sess.items[key]
	return
}

// 获取值，若键名对应的值不存在，则返回defVal。
func (sess *Session) MustGet(key, defVal interface{}) interface{} {
	sess.Lock()
	defer sess.Unlock()

	val, found := sess.items[key]
	if !found {
		return defVal
	}
	return val
}

// 添加或是设置值。
func (sess *Session) Set(key, val interface{}) {
	sess.Lock()
	defer sess.Unlock()

	sess.items[key] = val
}

// 指定的键值是否存在。
func (sess *Session) Exists(key interface{}) bool {
	sess.Lock()
	defer sess.Unlock()

	_, found := sess.items[key]
	return found
}

// 当前session的sessionid
func (sess *Session) ID() string {
	return sess.id
}

func (sess *Session) Free(w http.ResponseWriter, req *http.Request) error {
	if err := sess.options.store.Delete(sess); err != nil {
		return err
	}

	sess.options.setCookie(w, sess.ID(), -1)
	return nil
}

// 保存当前的Session值到Options.store中。
func (sess *Session) Save(w http.ResponseWriter, req *http.Request) error {
	return sess.options.store.Save(sess)
}

// 产生一个唯一的SessionID
func sessionID() (string, error) {
	ret := make([]byte, 64)
	n, err := io.ReadFull(rand.Reader, ret)
	if n == 0 {
		return "", errors.New("未读取到随机数")
	}

	h := md5.New()
	h.Write(ret)
	return hex.EncodeToString(h.Sum(nil)), err
}
