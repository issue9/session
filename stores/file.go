// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package stores

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

// session文件创建的权限。
const mode os.FileMode = 0600

type file struct {
	dir string // session保存的路径
}

// 声明一个实现session.Store接口的文件存储器，
// 在该存储器下，每个session都将以单独的文件存储。
// dir为session文件的存放路径。创建的文件权限默认为0600。
func NewFile(dir string) (*file, error) {
	stat, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(dir, mode)
		} else {
			return nil, err
		}

		// 无法创建该目录
		if stat, err = os.Stat(dir); err != nil && os.IsNotExist(err) {
			return nil, err
		}
	}

	if !stat.IsDir() {
		return nil, fmt.Errorf("%v存在，但不是一个有效的路径。", dir)
	}

	return &file{
		dir: dir + string(os.PathSeparator),
	}, nil
}

// 该文件是否不存在
func (f *file) isNotExists(path string) bool {
	_, err := os.Stat(path)

	if err != nil && os.IsNotExist(err) {
		return true
	}

	return false
}

// session.Store.Delete()
func (f *file) Delete(sessID string) error {
	path := f.dir + sessID

	if f.isNotExists(path) {
		return nil
	}

	return os.Remove(path)
}

// session.Store.Get()
func (f *file) Get(sessID string) (map[interface{}]interface{}, error) {
	path := f.dir + sessID

	if f.isNotExists(path) { // 不存在，返回一个空值
		return map[interface{}]interface{}{}, nil
	}

	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	d := gob.NewDecoder(fp)
	mapped := make(map[interface{}]interface{}, 0)
	if err := d.Decode(&mapped); err != nil {
		return nil, err
	}

	return mapped, nil
}

// session.Store.Save()
func (f *file) Save(sessID string, data map[interface{}]interface{}) error {
	path := f.dir + sessID

	context := new(bytes.Buffer)
	e := gob.NewEncoder(context)
	e.Encode(data)

	fp, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fp.Close()

	_, err = context.WriteTo(fp)
	return err
}

// session.Store.GC()
func (f *file) GC(maxAge int) error {
	d := time.Now().Add(-time.Second * time.Duration(maxAge))

	fs, err := ioutil.ReadDir(f.dir)
	if err != nil {
		return err
	}

	for _, info := range fs {
		if info.IsDir() {
			continue
		}

		if info.ModTime().After(d) { // 未过期： info.ModTime() > (time.Now() - maxAge)
			continue
		}

		// 过期
		if err := os.Remove(f.dir + info.Name()); err != nil {
			return err
		}
	}
	return nil
}

// session.Store.Free()
func (f *file) Free() error {
	fs, err := ioutil.ReadDir(f.dir)
	if err != nil {
		return err
	}

	for _, info := range fs {
		if info.IsDir() {
			continue
		}

		if err := os.Remove(f.dir + info.Name()); err != nil {
			return err
		}
	}

	return nil
}
