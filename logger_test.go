// SPDX-License-Identifier: ISC
// Copyright (c) 2014-2020 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package logger_test

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bitmark-inc/logger"
)

var testLevelMap = map[string]string{
	"main": "debug",
	"aux":  "warn",
}

const (
	logDirectory     = "log"
	logFileName      = "test.log"
	logSizeOfFiles   = 30000
	logNumberOfFiles = 10
)

func removeLogFiles() {
	pathName := path.Join(logDirectory, logFileName)
	os.Remove(pathName)
	for i := 0; i <= logNumberOfFiles; i += 1 {
		os.Remove(pathName + "." + strconv.Itoa(i))
	}
	os.Remove(logDirectory)
}

func setup(t *testing.T) {
	removeLogFiles()
	os.Mkdir(logDirectory, 0770)
	c := logger.Configuration{
		Directory: logDirectory,
		File:      logFileName,
		Size:      logSizeOfFiles,
		Count:     logNumberOfFiles,
		Levels:    testLevelMap,
	}

	err := logger.Initialise(c)
	if err != nil {
		t.Fatalf("Logger setup failed with error: %v", err)
	}
}

func teardown() {
	removeLogFiles()
}

func TestLevels(t *testing.T) {
	setup(t)
	defer teardown()

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
	defer teardown()

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
	pathName := path.Join(logDirectory, logFileName)
	f, err := os.Open(pathName)
	if err != nil {
		t.Errorf("Failed to open %s because: %v", pathName, err)
		return
	}
	defer f.Close()
	r := bufio.NewReader(f)

	// length of the data and time prefix to skip
	dateTimeLength := 19

loop:
	for _, line := range strings.Split(s, "\n") {
		actualLine, err := r.ReadString('\n')
		if err == io.EOF && line == "" {
			break loop
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

func TestListLevels(t *testing.T) {
	l, err := logger.ListLevels()
	assert.Nil(t, err, "wrong ListLevels")

	var s logger.LogLevels
	err = json.Unmarshal(l, &s)
	assert.Nil(t, err, "wrong bytes unmarshal")

	for _, l := range s.Levels {
		if l.Tag == "main" {
			assert.Equal(t, "debug", l.LogLevel, "wrong main level")
		}

		if l.Tag == "aux" {
			assert.Equal(t, "warn", l.LogLevel, "wrong aux level")
		}
	}
}

func TestUpdateTagLogLevel(t *testing.T) {
	mainLog := logger.New("main")
	mainLog.Debug("This should log")
	mainLog.Info("This should log")

	err := logger.UpdateTagLogLevel("main", "info")
	assert.Nil(t, err, "wrong UpdateTagLogLevel")

	mainLog.Trace("This should not log")
	mainLog.Debug("This should not log")
	mainLog.Info("This should log")
	mainLog.Warn("This should log")
	mainLog.Error("This should log")
	mainLog.Critical("This should log")

	// ensure panic won't cause test to fail
	defer func() {
		recover()
	}()

	logger.Panic("this should log")
}
