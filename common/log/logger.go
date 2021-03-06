/*
   This file is part of go-palletone.
   go-palletone is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.
   go-palletone is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.
   You should have received a copy of the GNU General Public License
   along with go-palletone.  If not, see <http://www.gnu.org/licenses/>.
*/
/*
 * @author PalletOne core developers <dev@pallet.one>
 * @date 2018
 */

// log is the palletone log system.

package log

import (
	"fmt"
	"github.com/palletone/go-palletone/common/files"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"strings"
	"sync"
	"time"
)

const (
	RootBuild      = "build"
	RootCmd        = "cmd"
	RootCommon     = "common"
	RootConfigure  = "configure"
	RootCore       = "core"
	RootInternal   = "internal"
	RootPtnclient  = "ptnclient"
	RootPtnjson    = "ptnjson"
	RootStatistics = "statistics"
	RootVendor     = "vendor"
	RootWallet     = "wallet"

	errorKey  = "ZAPLOG_ERROR"
	LogStdout = "stdout"
)

var defaultLogModule = []string{RootBuild, RootCmd, RootCommon, RootConfigure, RootCore, RootInternal, RootPtnclient, RootPtnjson, RootStatistics, RootVendor, RootWallet}

var Logger *zap.Logger
var mux sync.RWMutex

//var originFileName string

// init zap.logger
func InitLogger() {
	//date := fmt.Sprintf("%d-%d-%d", time.Now().Year(), time.Now().Month(), time.Now().Day())
	for i, path := range DefaultConfig.OutputPaths {
		if path == LogStdout {
			continue
		}
		//if originFileName == "" {
		//	originFileName = path
		//}
		//index := strings.LastIndex(originFileName, ".")
		//if -1 == index {
		//	index = len(originFileName)
		//}
		// DefaultConfig.OutputPaths[i] = fmt.Sprintf("%s_%s.%s", Substr(originFileName, 0, index), date, Substr(originFileName, index+1, len(originFileName)-index))

		if err := files.MakeDirAndFile(DefaultConfig.OutputPaths[i]); err != nil {
			panic(err)
		}
	}

	initLogger()
	log.SetFlags(log.Lmicroseconds | log.Lshortfile | log.LstdFlags)
	// go check()
}

func ConInitLogger() {

	DefaultConfig.LoggerLvl = "FATAL"
	initLogger()
	log.SetFlags(log.Lmicroseconds | log.Lshortfile | log.LstdFlags)
}

func FileInitLogger(logfile string) {
	DefaultConfig.OutputPaths = []string{logfile}
	initLogger()
	log.SetFlags(log.Lmicroseconds | log.Lshortfile | log.LstdFlags)
}

// init logger.
func initLogger() {
	var cfg zap.Config
	cfg.OutputPaths = DefaultConfig.OutputPaths
	cfg.ErrorOutputPaths = DefaultConfig.ErrorOutputPaths
	var lvl zap.AtomicLevel
	lvl.UnmarshalText([]byte(DefaultConfig.LoggerLvl))
	cfg.Level = lvl
	cfg.Encoding = DefaultConfig.Encoding
	cfg.Development = DefaultConfig.Development
	cfg.EncoderConfig = zap.NewProductionEncoderConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	//cfg.EncoderConfig.EncodeLevel=zapcore.LowercaseColorLevelEncoder
	l, err := cfg.Build()
	if err != nil {
		log.Fatal("init logger error: ", err)
	}
	// add openModule
	if strings.Contains(DefaultConfig.OpenModule[0], ",") {
		arr := strings.Split(DefaultConfig.OpenModule[0], ",")
		DefaultConfig.OpenModule[0] = ""
		DefaultConfig.OpenModule = append(DefaultConfig.OpenModule, arr...)
		DefaultConfig.OpenModule = append(DefaultConfig.OpenModule, defaultLogModule...)
	} else {
		if !(len(DefaultConfig.OpenModule) == 1 && DefaultConfig.OpenModule[0] == "all") {
			DefaultConfig.OpenModule = append(DefaultConfig.OpenModule, defaultLogModule...)
		}
	}
	l.SetOpenModule(DefaultConfig.OpenModule)
	l = l.WithOptions(zap.AddCallerSkip(1))
	if DefaultConfig.RotationMaxSize > 0 {
		includeStdout, filePath := getOutputPath(DefaultConfig.OutputPaths)
		rotateLogCore := func(core zapcore.Core) zapcore.Core {
			w := zapcore.AddSync(&lumberjack.Logger{
				Filename:   filePath,
				MaxSize:    DefaultConfig.RotationMaxSize, // megabytes
				MaxBackups: 3,
				MaxAge:     DefaultConfig.RotationMaxAge, // days
			})
			if includeStdout {
				stdout, _, _ := zap.Open("stdout")
				w = zap.CombineWriteSyncers(stdout, w)
			}
			encoder := zapcore.NewConsoleEncoder(cfg.EncoderConfig)
			core2 := zapcore.NewCore(
				encoder,
				w,
				cfg.Level,
			)
			return core2
		}
		l = l.WithOptions(zap.WrapCore(rotateLogCore))
	}
	Logger = l
}

func getOutputPath(paths []string) (bool, string) {
	includeStdout := false
	filePath := ""
	for _, path := range paths {
		if strings.ToLower(path) == "stdout" {
			includeStdout = true
		} else {
			filePath = path
		}
	}
	return includeStdout, filePath
}

// Trace
func Trace(msg string, ctx ...interface{}) {
	if Logger == nil {
		//log.Println("logger is nil.")
		InitLogger()
	}
	//log.Println("logger trace is  ok.")
	fileds := ctxTOfileds(ctx...)
	Logger.Debug(msg, fileds...)
}

// Debug
func Debug(msg string, ctx ...interface{}) {
	if Logger == nil {
		InitLogger()
	}
	fileds := ctxTOfileds(ctx...)
	Logger.Debug(msg, fileds...)
}
func Debugf(format string, ctx ...interface{}) {
	if Logger == nil {
		InitLogger()
	}
	Logger.Debug(fmt.Sprintf(format, ctx...))
}

// Info
func Info(msg string, ctx ...interface{}) {
	if Logger == nil {
		InitLogger()
	}
	fileds := ctxTOfileds(ctx...)
	Logger.Info(msg, fileds...)
}

func Infof(msg string, ctx ...interface{}) {
	if Logger == nil {
		InitLogger()
	}
	Logger.Info(fmt.Sprintf(msg, ctx...))
}

// Warn
func Warn(msg string, ctx ...interface{}) {
	if Logger == nil {
		InitLogger()
	}
	fileds := ctxTOfileds(ctx...)
	Logger.Warn(msg, fileds...)
}
func Warnf(msg string, ctx ...interface{}) {
	if Logger == nil {
		InitLogger()
	}
	Logger.Warn(fmt.Sprintf(msg, ctx...))
}

// Error
func Error(msg string, ctx ...interface{}) {
	if Logger == nil {
		InitLogger()
	}

	fileds := ctxTOfileds(ctx...)
	Logger.Error(msg, fileds...)
}

func Errorf(format string, ctx ...interface{}) {
	if Logger == nil {
		InitLogger()
	}
	Logger.Error(fmt.Sprintf(format, ctx...))
}

// Crit
func Crit(msg string, ctx ...interface{}) {
	if Logger == nil {
		InitLogger()
	}
	fileds := ctxTOfileds(ctx...)
	Logger.Error(msg, fileds...)
}

// ctx transfer to  fileds
func ctxTOfileds(ctx ...interface{}) []zap.Field {
	// ctx translate into zap.Filed
	normalctx := normalize(ctx)
	fileds := make([]zap.Field, 0)
	var prefix, suffix []interface{}
	for i, v := range normalctx {
		if i%2 == 0 {
			prefix = append(prefix, v)
		} else {
			suffix = append(suffix, v)
		}
	}

	for i := 0; i < len(prefix); i++ {
		pr := prefix[i]
		if e, ok := pr.(error); ok {
			fileds = append(fileds, zap.Any(e.Error(), suffix[i]))
		} else {
			fileds = append(fileds, zap.Any(pr.(string), suffix[i]))
		}
	}
	return fileds
}

// normalize
func normalize(ctx []interface{}) []interface{} {
	// if the caller passed a Ctx object, then expand it
	if len(ctx) == 1 {
		if ctxMap, ok := ctx[0].(Ctx); ok {
			ctx = ctxMap.toArray()
		}
	}

	// ctx needs to be even because it's a series of key/value pairs
	// no one wants to check for errors on logging functions,
	// so instead of erroring on bad input, we'll just make sure
	// that things are the right length and users can fix bugs
	// when they see the output looks wrong
	if len(ctx)%2 != 0 {
		ctx = append(ctx, nil, errorKey, "Normalized odd number of arguments by adding nil")
	}

	return ctx
}

// Lazy allows you to defer calculation of a logged value that is expensive
// to compute until it is certain that it must be evaluated with the given filters.
//
// Lazy may also be used in conjunction with a Logger's New() function
// to generate a child logger which always reports the current value of changing
// state.
//
// You may wrap any function which takes no arguments to Lazy. It may return any
// number of values of any type.
// type Lazy struct {
// 	Fn interface{}
// }

// Ctx is a map of key/value pairs to pass as context to a log function
// Use this only if you really need greater safety around the arguments you pass
// to the logging functions.
type Ctx map[string]interface{}

func (c Ctx) toArray() []interface{} {
	arr := make([]interface{}, len(c)*2)

	i := 0
	for k, v := range c {
		arr[i] = k
		arr[i+1] = v
		i += 2
	}

	return arr
}

func Substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}

func check() {
	for {
		time.Sleep(time.Duration(5) * time.Second)
		date := fmt.Sprintf("%d-%d-%d", time.Now().Year(), time.Now().Month(), time.Now().Day())
		mux.RLock()
		filename := Logger.GetFileName()
		mux.RUnlock()
		if strings.Index(filename, date) == -1 {
			mux.Lock()
			Logger = nil
			InitLogger()
			mux.Unlock()
			return
		}
	}
}

/*
type ILogger interface {
	Trace(msg string, ctx ...interface{})
	Debug(msg string, ctx ...interface{})
	Debugf(format string, ctx ...interface{})
	Info(msg string, ctx ...interface{})
	Infof(format string, ctx ...interface{})
	Warn(msg string, ctx ...interface{})
	Warnf(format string, ctx ...interface{})
	Error(msg string, ctx ...interface{})
	Errorf(format string, ctx ...interface{})
	Crit(msg string, ctx ...interface{})
}

type Plogger struct {
	logger zap.Logger
}

// New returns a new logger with the given context.
// New is a convenient alias for Root().New
func New(ctx ...interface{}) *Plogger {
	if Logger == nil {
		InitLogger()
	}
	pl := new(Plogger)
	pl.logger = *Logger
	return pl
}

func NewTestLog() *Plogger {
	DefaultConfig = Config{
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		OpenModule:       []string{"all"},
		LoggerLvl:        "DEBUG",
		Encoding:         "console",
		Development:      true,
	}
	initLogger()
	return &Plogger{logger: *Logger}
}
func (pl *Plogger) New(ctx ...interface{}) *Plogger {
	if pl != nil {
		return pl
	}
	if Logger == nil {
		InitLogger()
	}

	pl.logger = *Logger
	return pl
}
func (pl *Plogger) Trace(msg string, ctx ...interface{}) {
	fileds := ctxTOfileds(ctx...)
	pl.logger.Debug(msg, fileds...)
}

func (pl *Plogger) Debug(msg string, ctx ...interface{}) {
	fileds := ctxTOfileds(ctx...)
	pl.logger.Debug(msg, fileds...)
}
func (pl *Plogger) Debugf(format string, ctx ...interface{}) {
	pl.logger.Debug(fmt.Sprintf(format, ctx...))
}
func (pl *Plogger) Info(msg string, ctx ...interface{}) {
	fileds := ctxTOfileds(ctx...)
	pl.logger.Info(msg, fileds...)
}
func (pl *Plogger) Infof(format string, ctx ...interface{}) {
	pl.logger.Info(fmt.Sprintf(format, ctx...))
}
func (pl *Plogger) Warn(msg string, ctx ...interface{}) {
	fileds := ctxTOfileds(ctx...)
	pl.logger.Warn(msg, fileds...)
}
func (pl *Plogger) Warnf(format string, ctx ...interface{}) {
	pl.logger.Warn(fmt.Sprintf(format, ctx...))
}
func (pl *Plogger) Error(msg string, ctx ...interface{}) {
	fileds := ctxTOfileds(ctx...)
	pl.logger.Error(msg, fileds...)
}
func (pl *Plogger) Errorf(format string, ctx ...interface{}) {
	pl.logger.Error(fmt.Sprintf(format, ctx...))
}
func (pl *Plogger) Crit(msg string, ctx ...interface{}) {
	fileds := ctxTOfileds(ctx...)
	pl.logger.Error(msg, fileds...)
}
*/
