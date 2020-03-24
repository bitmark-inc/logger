// SPDX-License-Identifier: ISC
// Copyright (c) 2014-2020 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package level_test

import (
	"github.com/bitmark-inc/logger/level"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidLevels(t *testing.T) {
	assert.Equal(t, level.TraceLevel, level.ValidLevels[level.Trace], "wrong trace")
	assert.Equal(t, level.DebugLevel, level.ValidLevels[level.Debug], "wrong debug")
	assert.Equal(t, level.InfoLevel, level.ValidLevels[level.Info], "wrong info")
	assert.Equal(t, level.WarnLevel, level.ValidLevels[level.Warn], "wrong warn")
	assert.Equal(t, level.ErrorLevel, level.ValidLevels[level.Error], "wrong error")
	assert.Equal(t, level.CriticalLevel, level.ValidLevels[level.Critical], "wrong critical")
	assert.Equal(t, level.OffLevel, level.ValidLevels[level.Off], "wrong off")
}
