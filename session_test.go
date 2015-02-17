// Copyright 2014 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package session

import (
	"testing"

	"github.com/issue9/assert"
)

func TestSessionID(t *testing.T) {
	m := make(map[string]interface{}, 0)

	// 随机产生几个字符串，看是否有可能重复
	for i := 0; i < 10000; i++ {
		sid, err := sessionID()
		assert.Nil(t, err)
		//assert.Equal(t, len(sid), sessionIDLen)

		_, found := m[sid]
		assert.False(t, found)

		m[sid] = nil
	}
}
