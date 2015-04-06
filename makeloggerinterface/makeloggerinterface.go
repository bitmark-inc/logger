// Copyright (c) 2014-2015 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

// generate the various logging calls
package main

import (
	"os"
	"strings"
	"text/template"
	"time"
)

const (
	header = `// Copyright (c) 2014-{{.Year}} Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

// Generated file: Do _NOT_ Modify
// Generated on: {{.Today}}

package logger
`
)

const (
	codeBlock = `
// Log a simple string
// e.g.
//   log.{{.CapitalLevel}}("a log message")
func (l *LoggerChannel) {{.CapitalLevel}}(message string) {
	if l.levelNumber <= {{.LowerLevel}}Value {
		l.log.{{.CapitalLevel}}(l.formatPrefix + message)
	}
}

// Log a formatted string with arguments lige fmt.Sprintf()
// e.g.
//   log.{{.CapitalLevel}}f("the value = %d", xValue)
func (l *LoggerChannel) {{.CapitalLevel}}f(format string, arguments ...interface{}) {
	if l.levelNumber <= {{.LowerLevel}}Value {
		l.log.{{.CapitalLevel}}f(l.formatPrefix+format, arguments...)
	}
}

// Log from a closure, any function returning a string is suitable
// and the closure will only be executed if the log level is low enough.
// This allows complex
// e.g.
//   log.{{.CapitalLevel}}c(func() string {
//       return fmt.Sprintf("the sin(%f) = %f", x, math.sin(x))
//   })
func (l *LoggerChannel) {{.CapitalLevel}}c(closure func() string) {
	if l.levelNumber <= {{.LowerLevel}}Value {
		l.log.{{.CapitalLevel}}(l.formatPrefix + closure())
	}
}
`
)

var levels = []string{
	"Trace",
	"Debug",
	"Info",
	"Warn",
	"Error",
	"Critical",
}

type headerExpansion struct {
	Today string
	Year  string
}

type expansion struct {
	CapitalLevel string
	LowerLevel   string
}

func main() {
	now := time.Now()
	headerParameters := headerExpansion{
		Today: now.Format(time.RFC3339),
		Year:  now.Format("2006"),
	}

	ht, err := template.New("interface").Parse(header)
	if err != nil {
		panic(err)
	}
	err = ht.Execute(os.Stdout, headerParameters)
	if err != nil {
		panic(err)
	}

	for _, level := range levels {
		parameters := expansion{
			CapitalLevel: level,
			LowerLevel:   strings.ToLower(level),
		}

		t, err := template.New("interface").Parse(codeBlock)
		if err != nil {
			panic(err)
		}
		err = t.Execute(os.Stdout, parameters)
		if err != nil {
			panic(err)
		}
	}
}
