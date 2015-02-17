// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package session

type Store interface {
	Delete(sess *Session) error

	// 获取指定ID的值
	Get(sid string) (*Session, error)

	// 将Session中的值保存到当前的实例中
	Save(*Session) error
}
