// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// session的操作包。
//
// 当前session包的sessionid只能通过cookie传递，
// 所以只有在启用了cookie的浏览器上，本session包才有用。
//
// 用户可以通过实现Store接口，自行实现Session数据的存储，
// 具体的实现方式可以参考stores目录下的相关实例，
// 该包实现了一些常用的Store。
//
// 以下是一个简单的session操作示例：
//  mgr := session.New(stores.NewMemory(...), providers.NewCookie(...))
//
//  h := func(w http.ResponseWriter, req *http.Request) {
//      // 在每一个Handler中调用Start()开始一个Session操作。
//      sess,err :=mgr.Start( w, req)
//      defer sess.Close()
//
//      sess.Get(...)
//  }
//  http.HandleFunc("/", h)
//  http.ListenAndServe(":8080")
//
//  // 服务结束后，记得释放Options实例。
//  mgr.Close()
//
// 也可以多个store同时使用：
//  frontMgr := session.New(stores.NewMemory(), providers.NewCookie())
//  adminMgr := session.New(stores.NewFile(), provider.NewCookie())
//
//  frontHandler := func(w http.ResponseWriter, req *http.Request) {
//      sess,err :=frontMgr.Start(w, req)
//      defer sess.Close()
//
//      sess.Get(...)
//  }
//
//  adminHandler := func(w http.ResponseWriter, req *http.Request) {
//      sess,err :=adminMgr.Start(w, req)
//      defer sess.Close()
//
//      sess.Get(...)
//  }
//
//  http.HandleFunc("/front", frontHandler)
//  http.HandleFunc("/admin", adminHandler)
//  http.ListenAndServe(":88")
//
//  frontMgr.Close()
//  adminMgr.Close()
package session
