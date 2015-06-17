// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package providers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/session/types"
)

var _ types.Provider = &cookie{}

func newCookie(a *assert.Assertion) types.Provider {
	provider := NewCookie(11, "gosession", "/", "localhost", false)
	a.NotNil(provider)
	return provider
}

func testCookie_Init(t *testing.T) {
	a := assert.New(t)
	//provider := newCookie(a)

	// 该handler通过action参数确定是设置还是删除cookie
	setCookieHandler := func(w http.ResponseWriter, req *http.Request) {
		a.NotError(req.ParseForm())
		//maxAge := -1
		switch req.Form["action"][0] {
		case "set":
			//maxAge = 10
		case "unset":
			//maxAge = -1
		default:
			t.Errorf("无效的action值:%v", req.Form["action"][0])
		}
		//provider.setCookie(w, "sessionIDValue", maxAge)
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

// 测试Cookie.getSessionID，是否能根据request还原相应的sessionid。
func TestCookie_Get(t *testing.T) {
	a := assert.New(t)

	provider := newCookie(a)

	var sid string // 记录相关的sessionid值。
	var err error
	h := func(w http.ResponseWriter, req *http.Request) {
		a.NotError(req.ParseForm())
		switch req.Form["action"][0] {
		case "1": // 第一次访问，设置sessionID
			sid, err = provider.Get(w, req)
			a.NotError(err).NotEmpty(sid)
		case "2": // 第二次访问，验证sessionID
			sessID, err := provider.Get(w, req)
			a.NotError(err).Equal(sessID, sid, "action=2:sessID[%v] != sid[%v]", sessID, sid)
			sid = sessID
		case "3":
			provider.Delete(w, req)
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

	// ?action=3，第三次访问，清空cookie。
	resp, err = http.Get(srv.URL + "?action=3")
	a.NotError(err).NotNil(resp)
	a.True(strings.Index(resp.Header.Get("Set-Cookie"), "Max-Age=0") >= 0)
}
