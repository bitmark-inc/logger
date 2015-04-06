// Copyright (c) 2014-2015 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

// A simple logging layer on top of seelog
//
// Implements multiple log channels to a single rotated log file.  Log
// levels are considered a simple hierarchy with each channel having a
// single limit, below which logs to that channel are skipped.
package logger

import (
	"fmt"
	"github.com/cihub/seelog"
	"strings"
	"sync"
)

// The log messages will be "<tag><tagSuffix><message>"
const tagSuffix = ": "

// the tagname reserved to set the default level for unknown tags
const DefaultTag = "*"

// the initial level for unknown tags
const DefaultLevel = "error"

// Holds a set of default values for future NewChannel calls
var levelMap = map[string]string{DefaultTag: DefaultLevel}

// simple ordering to allow <= to decide if log will be output
const (
	_ = iota
	traceValue
	debugValue
	infoValue
	warnValue
	errorValue
	criticalValue
	offValue
)

// This needs to correspond to seelog levels
var validLevels = map[string]int{
	"trace":    traceValue,
	"debug":    debugValue,
	"info":     infoValue,
	"warn":     warnValue,
	"error":    errorValue,
	"critical": criticalValue,
	"off":      offValue,
}

// The logging structure
// nice short name so use like:
//
//   var log *logger.L
//   log := logger.New("sometag")
type L struct {
	sync.Mutex
	tag          string
	formatPrefix string
	textPrefix   string
	level        string
	levelNumber  int
	log          seelog.LoggerInterface
}

// Pre-load default levels before creating any new logging channels
// invalid levels are simply skipped and repeated calls will
// accumulate new tag values and overwrite old tag values.
//
// This will not update currently open channels, it is only for new
// channels, it is intended to be called once after command-line
// arguments and any configuration files have been processed to
// establish logging defaults.
//
// the name "*" is reserved to set a level for any tags that do not
// have table entries.
func LoadLevels(levels map[string]string) {
	for tag, level := range levels {
		// make sure that levelMap only contains correct data
		// by ignoring invalid levels
		if _, ok := validLevels[level]; ok {
			levelMap[tag] = level
		}
	}
}

// Setup seelog to used a rotated log file and output all logs
// level control is now controlled by this module
func Initialise(file string, size int, number int) error {
	config := fmt.Sprintf(`
          <seelog type="adaptive"
                  mininterval="2000000"
                  maxinterval="100000000"
                  critmsgcount="500"
                  minlevel="trace">
              <outputs formatid="all">
                  <rollingfile type="size" filename="%s" maxsize="%d" maxrolls="%d" />
              </outputs>
              <formats>
                  <format id="all" format="%%Date %%Time [%%LEVEL] %%Msg%%n" />
              </formats>
          </seelog>`, file, size, number)

	logger, err := seelog.LoggerFromConfigAsString(config)
	if err != nil {
		return err
	}
	err = seelog.ReplaceLogger(logger)
	if nil == err {
		seelog.Current.Warn("LOGGER: ===== Logging system started =====")
	}
	return err
}

// flush all channels
func Finalise() {
	seelog.Current.Warn("LOGGER: ===== Logging system stopped =====")
	seelog.Flush()
}

// flush all channels
func Flush() {
	seelog.Flush()
}

// Open a new logging channel with a specified tag
func New(tag string) *L {

	// map % -> %% to be printf safe
	s := strings.Split(tag, "%")
	j := strings.Join(s, "%%")

	// determine the initial level
	level, ok := levelMap[tag]
	if !ok {
		level, ok = levelMap[DefaultTag]
	}
	if !ok {
		level = "error"
	}

	// create a logger channel
	return &L{
		tag:          tag, // for referencing default level
		formatPrefix: j + tagSuffix,
		textPrefix:   tag + tagSuffix,
		level:        level,
		levelNumber:  validLevels[level], // level is validated so get a non-zero value
		log:          seelog.Current,
	}
}

// Change the log level for a given channel returns the current level
//
// Use the value of DefaultTag to return to current default value.
// Use the value "" to just return the current setting
func (l *L) ChangeLevel(level string) string {
	// preserve current
	current := l.level

	// to return to default level which may have been modified
	// by subsequent LoadLevels calls.
	if DefaultTag == level {
		// get currrent default level for this tag
		var ok bool
		level, ok = levelMap[l.tag]
		if !ok {
			level, ok = levelMap[DefaultTag]
		}
		if !ok {
			level = DefaultLevel
		}
	} else if "" == level {
		return current
	}

	// set level and corresponding number
	if n, ok := validLevels[level]; ok {
		l.Lock()
		defer l.Unlock()
		l.level = level
		l.levelNumber = n
	}
	return current
}

// flush messages
func (l *L) Flush() {
	Flush()
}
