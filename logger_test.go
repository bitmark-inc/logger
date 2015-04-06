// Copyright (c) 2014-2015 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package logger_test

import (
	"bufio"
	"github.com/bitmark-inc/logger"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"
)

var testLevelMap = map[string]string{
	"main": "debug",
	"aux":  "warn",
}

const (
	logFileName      = "test.log"
	logNumberOfFiles = 3
)

func removeLogFiles() {
	os.Remove(logFileName)
	for i := 0; i <= logNumberOfFiles; i += 1 {
		os.Remove(logFileName + "." + strconv.Itoa(i))
	}
}

func setup(t *testing.T) {
	removeLogFiles()
	err := logger.Initialise(logFileName, 1024, 2)
	if err != nil {
		t.Errorf("Logger setup failed with error: %v", err)
	}
	logger.LoadLevels(testLevelMap)
}

func teardown(t *testing.T) {
	removeLogFiles()
}

func TestLevels(t *testing.T) {
	setup(t)
	defer teardown(t)

	mainLog := logger.New("main")
	auxLog := logger.New("aux")

	mainLog.Trace("This should not log")
	mainLog.Debug("This should log")
	mainLog.Info("This should log")
	mainLog.Warn("This should log")
	mainLog.Error("This should log")
	mainLog.Critical("This should log")

	auxLog.Trace("This should not log")
	auxLog.Debug("This should not log")
	auxLog.Info("This should not log")
	auxLog.Warn("This should log")
	auxLog.Error("This should log")
	auxLog.Critical("This should log")

	checkfile(t, `2014-08-12 10:44:35 [WARN] LOGGER: ===== Logging system started =====
2014-08-12 10:44:35 [DEBUG] main: This should log
2014-08-12 10:44:35 [INFO] main: This should log
2014-08-12 10:44:35 [WARN] main: This should log
2014-08-12 10:44:35 [ERROR] main: This should log
2014-08-12 10:44:35 [CRITICAL] main: This should log
2014-08-12 10:44:35 [WARN] aux: This should log
2014-08-12 10:44:35 [ERROR] aux: This should log
2014-08-12 10:44:35 [CRITICAL] aux: This should log
2014-08-12 10:44:35 [WARN] LOGGER: ===== Logging system stopped =====
`)
}

func TestClosure(t *testing.T) {
	setup(t)
	defer teardown(t)

	mainLog := logger.New("main")

	mainClosureVar := false

	// ensure closure does not execute
	mainLog.Tracec(func() string {
		mainClosureVar = true
		return "This should not log"
	})

	if mainClosureVar {
		t.Errorf("closure was called - when it should not")
		return
	}

	// ensure closure does execute
	mainLog.Warnc(func() string {
		mainClosureVar = true
		return "This should log"
	})

	if !mainClosureVar {
		t.Errorf("closure was not called - when it should")
		return
	}

	checkfile(t, `2014-08-12 10:44:35 [WARN] LOGGER: ===== Logging system started =====
2014-08-12 10:44:35 [WARN] main: This should log
2014-08-12 10:44:35 [WARN] LOGGER: ===== Logging system stopped =====
`)
}

// compare actual log results with expected, ignoring the dat and time values
func checkfile(t *testing.T, s string) {
	logger.Finalise()
	f, err := os.Open(logFileName)
	if err != nil {
		t.Errorf("Failed to open %s because: %v", logFileName, err)
		return
	}
	defer f.Close()
	r := bufio.NewReader(f)

	// length of the data and time prefix to skip
	dateTimeLength := 19

	for _, line := range strings.Split(s, "\n") {
		actualLine, err := r.ReadString('\n')
		if err == io.EOF && line == "" {
			break
		}
		if err != nil {
			t.Errorf("Error reading %s : %v", logFileName, err)
			return
		}
		actualLine = actualLine[dateTimeLength : len(actualLine)-1] // trim '\n'
		if actualLine != line[dateTimeLength:] {
			t.Errorf("Mismatch read: '%s' wanted: '%s'", actualLine, line)
		}
	}
}
