// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package session

// Session的存储接口。
type Store interface {
	// 删除指定的Session，若是存在的话。
	// 若Session不存在，则不发生任何事情。
	Delete(sess *Session) error

	// 获取指定ID的Session实例。若不存在，则创建一个新的。
	Get(sid string) (*Session, error)

	// 将Session中的值保存到当前的实例中
	Save(*Session) error
}
