// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package options

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
)

// 产生一个唯一值。
func sessionID() (string, error) {
	ret := make([]byte, 64)
	n, err := io.ReadFull(rand.Reader, ret)
	if n == 0 {
		return "", errors.New("未读取到随机数")
	}

	h := md5.New()
	h.Write(ret)
	return hex.EncodeToString(h.Sum(nil)), err
}
