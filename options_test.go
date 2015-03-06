// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package session

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/issue9/assert"
)

func newOptions(a *assert.Assertion) *Options {
	opt := NewOptions(newTestStore(), 10, "gosession", "/", "localhost", false)
	a.NotNil(opt)
	return opt
}

func TestOptions_setCookie(t *testing.T) {
	a := assert.New(t)

	opt := newOptions(a)
	defer opt.Close()

	// 该handler通过action参数确定是设置还是删除cookie
	setCookieHandler := func(w http.ResponseWriter, req *http.Request) {
		a.NotError(req.ParseForm())
		maxAge := -1
		switch req.Form["action"][0] {
		case "set":
			maxAge = 10
		case "unset":
			maxAge = -1
		default:
			t.Errorf("无效的action值:%v", req.Form["action"][0])
		}
		opt.setCookie(w, "sessionIDValue", maxAge)
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	}

	srv := httptest.NewServer(http.HandlerFunc(setCookieHandler))
	a.NotNil(srv)
	defer srv.Close()

	/* ?action=set，设置cookie值 */

	response, err := http.Get(srv.URL + "?action=set")
	a.NotError(err).NotNil(response)

	var cookie *http.Cookie = nil
	cookies := response.Cookies()
	for _, c := range cookies {
		if c.Name == "gosession" {
			cookie = c
			break
		}
	}
	a.NotNil(cookie, "不存在gosession的cookie，response.Cookies()值:%v", cookies)
	a.True(cookie.MaxAge > 0, "cookie.MaxAge>0不成立；cookie的值为:%v", cookie)

	/* ?action=unset，取消cookie值 */

	// 构建带cookie值的request
	r, err := http.NewRequest("GET", srv.URL+"?action=unset", nil)
	a.NotError(err).NotNil(r)
	r.Header.Add("Cookie", response.Header.Get("Set-Cookie"))
	client := &http.Client{}

	response, err = client.Do(r)
	a.NotError(err).NotNil(response)
	srv.Close()

	cookie = nil
	cookies = response.Cookies()
	for _, c := range cookies {
		if c.Name == "gosession" {
			cookie = c
			break
		}
	}
	a.NotNil(cookie, "依然存在cookie值:%v", cookie)
	a.True(cookie.MaxAge == -1, "cookie.MaxAge==-1不成立,cookie.MaxAge的值为:%v", cookie.MaxAge)
}

// 测试Options.getSessionID，是否能根据request还原相应的sessionid。
func TestOptions_getSessionID1(t *testing.T) {
	a := assert.New(t)

	opt := newOptions(a)
	defer opt.Close()

	var sid string // 记录相关的sessionid值。
	var err error
	h := func(w http.ResponseWriter, req *http.Request) {
		a.NotError(req.ParseForm())
		switch req.Form["action"][0] {
		case "1": // 第一次访问，设置sessionID
			sid, err = opt.getSessionID(req)
			a.NotError(err).NotEmpty(sid)

			opt.setCookie(w, sid, 10)
		case "2": // 第二次访问，验证sessionID
			sessID, err := opt.getSessionID(req)
			a.NotError(err).Equal(sessID, sid)
		default:
			t.Errorf("无效的action值:%v", req.Form["action"][0])
		}
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	}
	srv := httptest.NewServer(http.HandlerFunc(h))
	a.NotNil(srv)
	defer srv.Close()

	// ?action=1，第一次访问，设置相应的session值。
	resp, err := http.Get(srv.URL + "?action=1")
	a.NotError(err).NotNil(resp)

	// ?action=2，第二次访问，带上上次返回的cookie值
	// 构建带cookie值的request
	r, err := http.NewRequest("GET", srv.URL+"?action=2", nil)
	a.NotError(err).NotNil(r)
	r.Header.Add("Cookie", resp.Header.Get("Set-Cookie"))
	client := &http.Client{}
	_, err = client.Do(r)
	a.NotError(err)
}

// 测试Options.getSessionID函数中随机数产生是否正常。
func TestOptions_getSessionID2(t *testing.T) {
	a := assert.New(t)

	opt := newOptions(a)
	defer opt.Close()

	req, err := http.NewRequest("GET", "/", nil)
	a.NotError(err).NotNil(req)

	m := make(map[string]interface{}, 0)

	// 随机产生几个字符串，看是否有可能重复
	for i := 0; i < 10000; i++ {
		sid, err := opt.getSessionID(req)
		a.Nil(err)

		_, found := m[sid]
		a.False(found)

		m[sid] = nil
	}
}

// 测试opt是否会正确执行sotre.GC()。
func TestOptions_GC(t *testing.T) {
	a := assert.New(t)

	store := newTestStore()
	opt := NewOptions(store, 3, "gosession", "/", "localhost", false)
	defer opt.Close()

	h := func(w http.ResponseWriter, req *http.Request) {
		a.NotError(req.ParseForm())
		switch req.Form["action"][0] {
		case "1": // 第一次访问
			sess, err := Start(opt, w, req)
			a.NotError(err).NotNil(sess)
			a.Equal(0, len(store.items)) // 保存之前为长度0
			sess.Save(w, req)
			a.Equal(1, len(store.items)) // 保存之后长度为1
		case "2": // 第二次访问，未超时。
			a.Equal(1, len(store.items))
		case "3": // 第二次访问，超时了，store被清空。
			a.Equal(0, len(store.items))
		default:
			t.Errorf("无效的action值:%v", req.Form["action"][0])
		}
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	}
	srv := httptest.NewServer(http.HandlerFunc(h))
	a.NotNil(srv)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "?action=1")
	a.NotError(err).NotNil(resp)
	resp, err = http.Get(srv.URL + "?action=2")
	a.NotError(err).NotNil(resp)
	time.Sleep(time.Second * 3) // 延时3秒
	resp, err = http.Get(srv.URL + "?action=3")
	a.NotError(err).NotNil(resp)
}
