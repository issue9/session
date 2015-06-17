// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package types

import (
	"net/http"
)

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

	// 启用GC。
	StartGC()

	// 释放整个Store存储的内容及关闭所有的GC操作。
	// 之后对Store的操作都将是未定义的。
	Close() error
}

// 提供sessionid的传递和保管。
// 一般为通过cookie传递，也有可能是通过url参数传递。
type Provider interface {
	// 从req中获取sessionid的值。或当sessionid不存在时，产生一个新值。
	Get(w http.ResponseWriter, req *http.Request) (sessID string, err error)

	// 删除当前保存的sessionid值。
	Delete(w http.ResponseWriter, req *http.Request) error
}
