// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package options

import (
	"testing"

	"github.com/issue9/assert"
)

func TestSessionID(t *testing.T) {
	a := assert.New(t)

	m := make(map[string]interface{}, 0)

	// 随机产生几个字符串，看是否有可能重复
	for i := 0; i < 10000; i++ {
		sid, err := sessionID()
		a.Nil(err)

		_, found := m[sid]
		a.False(found)

		m[sid] = nil
	}
}
