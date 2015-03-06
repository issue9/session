// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package session

// Session的存储接口。
//
// 不能将一个Store实例与多个Options实例进行关联。
// 否则可能会造成Session数据相互覆盖的情况。
type Store interface {
	// 从Store中删除指定sessionid的数据。
	Delete(sessID string) error

	// 获取与sessID关联的数据，若不存在，则返回空的map值。
	Get(sessID string) (map[interface{}]interface{}, error)

	// 将data与sessID相关联，并保存到当前Store实例中。
	Save(sessID string, data map[interface{}]interface{}) error

	// 回收超过时间的数据。
	GC(maxAge int) error

	// 释放整个Store存储的内容，之后对Store的操作都将是未定义的。
	Free() error
}
