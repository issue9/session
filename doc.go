// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// session的操作包。
//
//  opt := session.NewOptions(stores.NewMemory(), ...)
//
//  // 在每一个Handler中调用Start()开始一个Session操作。
//  sess,err :=session.Start(opt, w, req)
//  defer sess.Close()
//
//  sess.Get(...)
//
//  // 关闭opt
//  opt.Close()
//
// 通过实现Store接口，可以实现自定义的存储系统。
package session

// 当前库的版本号
const Version = "0.5.4.150306"
