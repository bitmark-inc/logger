// Copyright (c) 2014-2017 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

// Generated file: Do _NOT_ Modify
// Generated on: 2017-09-19T11:23:21+08:00

package logger

import (
	"fmt"
)


// Log a simple string
// e.g.
//   log.Trace("a log message")
func (l *L) Trace(message string) {
	if !logInitialised || nil == l {
		panic("logger is not initialised")
	}
	if l.levelNumber <= traceValue {
		l.log.Trace(l.formatPrefix + message)
	}
}

// Log a formatted string with arguments lige fmt.Sprintf()
// e.g.
//   log.Tracef("the value = %d", xValue)
func (l *L) Tracef(format string, arguments ...interface{}) {
	if !logInitialised || nil == l {
		panic("logger is not initialised")
	}
	if l.levelNumber <= traceValue {
		s := fmt.Sprintf(l.formatPrefix+format, arguments...)
		l.log.Trace(s)
	}
}

// Log from a closure, any function returning a string is suitable
// and the closure will only be executed if the log level is low enough.
// This allows complex
// e.g.
//   log.Tracec(func() string {
//       return fmt.Sprintf("the sin(%f) = %f", x, math.sin(x))
//   })
func (l *L) Tracec(closure func() string) {
	if !logInitialised || nil == l {
		panic("logger is not initialised")
	}
	if l.levelNumber <= traceValue {
		l.log.Trace(l.formatPrefix + closure())
	}
}

// Log a simple string
// e.g.
//   log.Debug("a log message")
func (l *L) Debug(message string) {
	if !logInitialised || nil == l {
		panic("logger is not initialised")
	}
	if l.levelNumber <= debugValue {
		l.log.Debug(l.formatPrefix + message)
	}
}

// Log a formatted string with arguments lige fmt.Sprintf()
// e.g.
//   log.Debugf("the value = %d", xValue)
func (l *L) Debugf(format string, arguments ...interface{}) {
	if !logInitialised || nil == l {
		panic("logger is not initialised")
	}
	if l.levelNumber <= debugValue {
		s := fmt.Sprintf(l.formatPrefix+format, arguments...)
		l.log.Debug(s)
	}
}

// Log from a closure, any function returning a string is suitable
// and the closure will only be executed if the log level is low enough.
// This allows complex
// e.g.
//   log.Debugc(func() string {
//       return fmt.Sprintf("the sin(%f) = %f", x, math.sin(x))
//   })
func (l *L) Debugc(closure func() string) {
	if !logInitialised || nil == l {
		panic("logger is not initialised")
	}
	if l.levelNumber <= debugValue {
		l.log.Debug(l.formatPrefix + closure())
	}
}

// Log a simple string
// e.g.
//   log.Info("a log message")
func (l *L) Info(message string) {
	if !logInitialised || nil == l {
		panic("logger is not initialised")
	}
	if l.levelNumber <= infoValue {
		l.log.Info(l.formatPrefix + message)
	}
}

// Log a formatted string with arguments lige fmt.Sprintf()
// e.g.
//   log.Infof("the value = %d", xValue)
func (l *L) Infof(format string, arguments ...interface{}) {
	if !logInitialised || nil == l {
		panic("logger is not initialised")
	}
	if l.levelNumber <= infoValue {
		s := fmt.Sprintf(l.formatPrefix+format, arguments...)
		l.log.Info(s)
	}
}

// Log from a closure, any function returning a string is suitable
// and the closure will only be executed if the log level is low enough.
// This allows complex
// e.g.
//   log.Infoc(func() string {
//       return fmt.Sprintf("the sin(%f) = %f", x, math.sin(x))
//   })
func (l *L) Infoc(closure func() string) {
	if !logInitialised || nil == l {
		panic("logger is not initialised")
	}
	if l.levelNumber <= infoValue {
		l.log.Info(l.formatPrefix + closure())
	}
}

// Log a simple string
// e.g.
//   log.Warn("a log message")
func (l *L) Warn(message string) {
	if !logInitialised || nil == l {
		panic("logger is not initialised")
	}
	if l.levelNumber <= warnValue {
		l.log.Warn(l.formatPrefix + message)
	}
}

// Log a formatted string with arguments lige fmt.Sprintf()
// e.g.
//   log.Warnf("the value = %d", xValue)
func (l *L) Warnf(format string, arguments ...interface{}) {
	if !logInitialised || nil == l {
		panic("logger is not initialised")
	}
	if l.levelNumber <= warnValue {
		s := fmt.Sprintf(l.formatPrefix+format, arguments...)
		l.log.Warn(s)
	}
}

// Log from a closure, any function returning a string is suitable
// and the closure will only be executed if the log level is low enough.
// This allows complex
// e.g.
//   log.Warnc(func() string {
//       return fmt.Sprintf("the sin(%f) = %f", x, math.sin(x))
//   })
func (l *L) Warnc(closure func() string) {
	if !logInitialised || nil == l {
		panic("logger is not initialised")
	}
	if l.levelNumber <= warnValue {
		l.log.Warn(l.formatPrefix + closure())
	}
}

// Log a simple string
// e.g.
//   log.Error("a log message")
func (l *L) Error(message string) {
	if !logInitialised || nil == l {
		panic("logger is not initialised")
	}
	if l.levelNumber <= errorValue {
		l.log.Error(l.formatPrefix + message)
	}
}

// Log a formatted string with arguments lige fmt.Sprintf()
// e.g.
//   log.Errorf("the value = %d", xValue)
func (l *L) Errorf(format string, arguments ...interface{}) {
	if !logInitialised || nil == l {
		panic("logger is not initialised")
	}
	if l.levelNumber <= errorValue {
		s := fmt.Sprintf(l.formatPrefix+format, arguments...)
		l.log.Error(s)
	}
}

// Log from a closure, any function returning a string is suitable
// and the closure will only be executed if the log level is low enough.
// This allows complex
// e.g.
//   log.Errorc(func() string {
//       return fmt.Sprintf("the sin(%f) = %f", x, math.sin(x))
//   })
func (l *L) Errorc(closure func() string) {
	if !logInitialised || nil == l {
		panic("logger is not initialised")
	}
	if l.levelNumber <= errorValue {
		l.log.Error(l.formatPrefix + closure())
	}
}

// Log a simple string
// e.g.
//   log.Critical("a log message")
func (l *L) Critical(message string) {
	if !logInitialised || nil == l {
		panic("logger is not initialised")
	}
	if l.levelNumber <= criticalValue {
		l.log.Critical(l.formatPrefix + message)
	}
}

// Log a formatted string with arguments lige fmt.Sprintf()
// e.g.
//   log.Criticalf("the value = %d", xValue)
func (l *L) Criticalf(format string, arguments ...interface{}) {
	if !logInitialised || nil == l {
		panic("logger is not initialised")
	}
	if l.levelNumber <= criticalValue {
		s := fmt.Sprintf(l.formatPrefix+format, arguments...)
		l.log.Critical(s)
	}
}

// Log from a closure, any function returning a string is suitable
// and the closure will only be executed if the log level is low enough.
// This allows complex
// e.g.
//   log.Criticalc(func() string {
//       return fmt.Sprintf("the sin(%f) = %f", x, math.sin(x))
//   })
func (l *L) Criticalc(closure func() string) {
	if !logInitialised || nil == l {
		panic("logger is not initialised")
	}
	if l.levelNumber <= criticalValue {
		l.log.Critical(l.formatPrefix + closure())
	}
}
