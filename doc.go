// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// session的操作包。
//
//  memStore := stores.NewMemory()
//  opt := session.NewOptions(memStore,...)
//
//  // 在每一个Handler中调用Start()开始一个Session操作。
//  sess,err :=session.Start(opt, w, req)
//
// 通过实现Store接口，可以实现自定义的存储系统。
package session

// 当前库的版本号
const Version = "0.3.2.150219"
