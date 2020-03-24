// SPDX-License-Identifier: ISC
// Copyright (c) 2014-2020 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package level

// simple ordering to allow <= to decide if log will be outputTarget
const (
	_ = iota
	TraceLevel
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	CriticalLevel
	OffLevel
)

const (
	Trace    = "trace"
	Debug    = "debug"
	Info     = "info"
	Warn     = "warn"
	Error    = "error"
	Critical = "critical"
	Off      = "off"
)

// This needs to correspond to seelog levels
var ValidLevels = map[string]int{
	Trace:    TraceLevel,
	Debug:    DebugLevel,
	Info:     InfoLevel,
	Warn:     WarnLevel,
	Error:    ErrorLevel,
	Critical: CriticalLevel,
	Off:      OffLevel,
}
