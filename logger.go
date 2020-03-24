// Copyright (c) 2014-2018 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package logger

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bitmark-inc/logger/level"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/cihub/seelog"
)

// initial configuration for the logger
//
// example of use (with structure tags for config file parsing)
//   type AppConfiguration struct {
//     …
//     Logging logger.Configuration `libucl:"logging" hcl:"logging" json:"logging"`
//     …
//   }
//
//   err := logger.Initialise(conf.Logging)
//   if nil != err {  // if failed then display message and exit
//     exitwithstatus.Message("logger error: %s", err)
//   }
//   defer logger.Finalise()
//
// example of ucl/hcl configuration section
//   logging {
//     directory = "/var/lib/app/log"
//     file = "app.log"
//     size = 1048576
//     count = 50
//     #console = true # to duplicate messages to console (default false)
//     levels {
//       DEFAULT = "info"
//       system = "error"
//       main = "warn"
//     }
//   }
type Configuration struct {
	Directory string            `libucl:"directory" hcl:"directory" json:"directory"`
	File      string            `libucl:"file" hcl:"file" json:"file"`
	Size      int               `libucl:"size" hcl:"size" json:"size"`
	Count     int               `libucl:"count" hcl:"count" json:"count"`
	Levels    map[string]string `libucl:"levels" hcl:"levels" json:"levels"`
	Console   bool              `libucl:"console" hcl:"console" json:"console"`
}

// some restrictions on sizes
const (
	minimumSize  = 20000
	minimumCount = 10

	// The log messages will be "<tag><tagSuffix><message>"
	tagSuffix = ": "

	// the tagname reserved to set the default level for unknown tags
	DefaultTag = "DEFAULT"

	// the initial level for unknown tags
	DefaultLevel = "error"
)

// Holds a set of default values for future NewChannel calls
var levelMap = map[string]string{DefaultTag: DefaultLevel}

// The logging channel structure
// example of use
//
//   var log *logger.L
//   log := logger.New("sometag")
//
//   log.Debugf("value: %d", value)
type L struct {
	sync.Mutex
	tag          string
	formatPrefix string
	textPrefix   string
	level        string
	levelNumber  int
	log          seelog.LoggerInterface
}

// LogLevels - log levels
type LogLevels struct {
	Levels []Level `json:"levels"`
}

// Level - log level info
type Level struct {
	Tag      string `json:"tag"`
	LogLevel string `json:"category"`
}

type outputTarget int

const (
	stdOut outputTarget = iota + 1
	fileOut
)

// loggers - holds all logger info
type loggers struct {
	sync.Mutex
	globalLog   *L
	initialised bool
	data        []*L
	output      outputTarget
}

var globalData loggers

// default set output to standard out
func init() {
	globalData.output = stdOut
	stdLogger := defaultLogger()

	_ = seelog.ReplaceLogger(stdLogger)
	// ensure that the global critical/panic functions always written
	globalData.globalLog = New("PANIC")
	globalData.globalLog.level = level.Critical
	globalData.globalLog.levelNumber = level.ValidLevels[level.Critical]
}

func defaultLogger() seelog.LoggerInterface {
	config := fmt.Sprintf(`
          <seelog type="adaptive"
                  mininterval="2000000"
                  maxinterval="100000000"
                  critmsgcount="500"
                  minlevel="trace">
              <outputs formatid="all">
			      <console />
              </outputs>
              <formats>
                  <format id="all" format="%%Date %%Time [%%LEVEL] %%Msg%%n" />
              </formats>
          </seelog>`)

	logger, _ := seelog.LoggerFromConfigAsString(config)
	return logger
}

// Set up the logging system
func Initialise(configuration Configuration) error {
	if globalData.initialised {
		return errors.New("logger is already initialised")
	}

	globalData.output = fileOut

	if "" == configuration.Directory {
		return errors.New("Directory cannot be empty")
	}

	if "" == configuration.File {
		return errors.New("File cannot be empty")
	}

	d, f := path.Split(configuration.File)
	if "" != d && f != configuration.File {
		return fmt.Errorf("File: %q cannot be a path name", configuration.File)
	}

	if configuration.Size < minimumSize {
		return fmt.Errorf("Size: %d cannot be less than: %d", configuration.Size, minimumSize)
	}

	if configuration.Count < minimumCount {
		return fmt.Errorf("Count: %d cannot be less than: %d", configuration.Count, minimumCount)
	}

	info, err := os.Lstat(configuration.Directory)
	if nil != err {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("Directory: %q does not exist", configuration.Directory)
	}
	//permission := info.Mode().Perm()

	lockFile := path.Join(configuration.Directory, "log.test")
	os.Remove(lockFile)
	fd, err := os.Create(lockFile)
	if nil != err {
		return err
	}
	defer fd.Close()
	defer os.Remove(lockFile)
	n, err := fd.Write([]byte("0123456789"))
	if nil != err {
		return err
	}
	if 10 != n {
		return errors.New("unable to write to logging files")
	}

	for tag, l := range configuration.Levels {
		// make sure that levelMap only contains correct data
		// by ignoring invalid levels
		if _, ok := level.ValidLevels[l]; ok {
			levelMap[tag] = l
		}
	}

	optionalConsole := ""
	if configuration.Console {
		optionalConsole = "<console />"
	}

	filepath := path.Join(configuration.Directory, configuration.File)
	config := fmt.Sprintf(`
          <seelog type="adaptive"
                  mininterval="2000000"
                  maxinterval="100000000"
                  critmsgcount="500"
                  minlevel="trace">
              <outputs formatid="all">
                  <rollingfile type="size" filename="%s" maxsize="%d" maxrolls="%d" />
                  %s
              </outputs>
              <formats>
                  <format id="all" format="%%Date %%Time [%%LEVEL] %%Msg%%n" />
              </formats>
          </seelog>`, filepath, configuration.Size, configuration.Count, optionalConsole)

	logger, err := seelog.LoggerFromConfigAsString(config)
	if err != nil {
		return err
	}
	err = seelog.ReplaceLogger(logger)
	if nil == err {
		seelog.Current.Warn("LOGGER: ===== Logging system started =====")
		globalData.initialised = true

		// ensure that the global critical/panic functions always write to the log file
		globalData.globalLog = New("PANIC")
		globalData.globalLog.level = level.Critical
		globalData.globalLog.levelNumber = level.ValidLevels[level.Critical]
	}
	return err
}

// flush all channels and let log message goes to standard out
func Finalise() {
	_ = seelog.Current.Warn("LOGGER: ===== Logging system stopped =====")
	seelog.Flush()

	// if log message goes to file, make it back to standard output
	if globalData.output == fileOut {
		globalData.initialised = false
		globalData.output = stdOut
		_ = seelog.ReplaceLogger(defaultLogger())

		globalData.globalLog = New("PANIC")
		globalData.globalLog.level = level.Critical
		globalData.globalLog.levelNumber = level.ValidLevels[level.Critical]
	}

	globalData.data = globalData.data[:0]
}

// flush all channels
func Flush() {
	seelog.Flush()
}

// Open a new logging channel with a specified tag
func New(tag string) *L {
	// if log message goes to file, then Initialise should be called
	if globalData.output == fileOut && !globalData.initialised {
		panic("logger.New Initialise was not called")
	}

	// map % -> %% to be printf safe
	s := strings.Split(tag, "%")
	j := strings.Join(s, "%%")

	// determine the level
	l, ok := levelMap[tag]
	if !ok {
		l, ok = levelMap[DefaultTag]
	}
	if !ok {
		l = DefaultLevel
	}

	// create a logger channel
	ptr := &L{
		tag:          tag, // for referencing default level
		formatPrefix: j + tagSuffix,
		textPrefix:   tag + tagSuffix,
		level:        l,
		levelNumber:  level.ValidLevels[l], // level is validated so get a non-zero value
		log:          seelog.Current,
	}

	globalData.Lock()
	globalData.data = append(globalData.data, ptr)
	defer globalData.Unlock()

	return ptr
}

// flush messages
func (l *L) Flush() {
	Flush()
}

// global logging message
func Critical(message string) {
	globalData.globalLog.Critical(message)
}

// global logging formatted message
func Criticalf(format string, arguments ...interface{}) {
	globalData.globalLog.Criticalf(format, arguments...)
}

// global logging message + panic
func Panic(message string) {
	globalData.globalLog.Critical(message)
	Flush()
	time.Sleep(100 * time.Millisecond) // to allow logging outputTarget
	panic(message)
}

// global logging formatted message + panic
func Panicf(format string, arguments ...interface{}) {
	globalData.globalLog.Criticalf(format, arguments...)
	Flush()
	time.Sleep(100 * time.Millisecond) // to allow logging outputTarget
	panic(fmt.Sprintf(format, arguments...))
}

// conditional panic
func PanicIfError(message string, err error) {
	if nil == err {
		return
	}
	Panicf("%s failed with error: %v", message, err)
}

// ListLevels - return log level info in json format
// it's not lock protected, because it should be low frequency to list log
// levels
func ListLevels() ([]byte, error) {
	levels := make([]Level, 0)

	for _, l := range globalData.data {
		levels = append(levels, Level{
			Tag:      l.tag,
			LogLevel: l.level,
		})
	}

	ll := LogLevels{Levels: levels}
	bs, err := json.Marshal(ll)
	if nil != err {
		return []byte{}, err
	}

	return bs, nil
}

// UpdateTagLogLevel - update log level for specific tag
func UpdateTagLogLevel(tag, newLevel string) error {
	globalData.Lock()
	defer globalData.Unlock()

	for i, l := range globalData.data {
		if l.tag == tag {
			if num, ok := level.ValidLevels[newLevel]; !ok {
				return fmt.Errorf("level %s invalid", newLevel)
			} else {
				globalData.data[i].levelNumber = num
				globalData.data[i].level = newLevel
				return nil
			}
		}
	}

	return fmt.Errorf("tag %s not found", tag)
}

func validLogger(l *L) bool {
	if globalData.output == stdOut {
		return true
	}

	if !globalData.initialised || nil == l {
		return false
	}

	return true
}
