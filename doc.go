// Copyright (c) 2014-2017 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

// A simple logging layer on top of seelog
//
// Implements multiple log channels to a single rotated log file.  Log
// levels are considered a simple hierarchy with each channel having a
// single limit, below which logs to that channel are skipped.
package logger
