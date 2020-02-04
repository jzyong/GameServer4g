package log

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

//TODO 待完善

const (
	//日志级别
	debugLevel = 0
	infoLevel  = 1
	warnLevel  = 2
	errorLevel = 3
	fatalLevel = 4
	//日志级别
	//级别输出字符标示
	debugLevelStr = "[debug]"
	infoLevelStr  = "[info	]"
	warnLevelStr  = "[warn	]"
	errorLevelStr = "[error]"
	fatalLevelStr = "[fatal]"
)

//日志
type Logger struct {
	level  int
	logger *log.Logger
	file   *os.File
}

//新建
func New(logLevel string, pathName string, logName string, flag int) (*Logger, error) {
	var level int
	switch strings.ToLower(logLevel) {
	case "debug":
		level = debugLevel
	case "info":
		level = infoLevel
	case "warn":
		level = warnLevel
	case "error":
		level = errorLevel
	case "fatal":
		level = fatalLevel
	default:
		return nil, errors.New("unknown level:" + logLevel)
	}

	var innerLogger *log.Logger
	var file *os.File
	if logName != "" {
		now := time.Now()
		fileName := fmt.Sprintf("%v-%d%02d%02d_%02d_%02d_%02d.log", logName,
			now.Year(),
			now.Month(),
			now.Day(),
			now.Hour(),
			now.Minute(),
			now.Second())
		f, err := os.Create(fileName)
		if err != nil {
			return nil, err
		}
		innerLogger = log.New(file, "", flag)
		file = f
	} else {
		innerLogger = log.New(os.Stdout, "", flag)
	}

	logger := new(Logger)
	logger.level = level
	logger.logger = innerLogger
	logger.file = file
	return logger, nil
}

// It's dangerous to call the method on logging
func (logger *Logger) Close() {
	if logger.file != nil {
		logger.file.Close()
	}
	logger.logger = nil
	logger.file = nil
}

//内部调用输出
func (logger *Logger) doPrintf(level int, printLevel string, format string, a ...interface{}) {
	if level < logger.level {
		return
	}
	if logger.logger == nil {
		panic("logger closed")
	}

	format = printLevel + format
	logger.logger.Output(3, fmt.Sprintf(format, a...))

	if level == fatalLevel {
		os.Exit(1)
	}
}

//调试
func (logger *Logger) Debug(format string, a ...interface{}) {
	logger.doPrintf(debugLevel, debugLevelStr, format, a...)
}

//打印信息
func (logger *Logger) Info(format string, a ...interface{}) {
	logger.doPrintf(infoLevel, infoLevelStr, format, a...)
}

//警告
func (logger *Logger) Warn(format string, a ...interface{}) {
	logger.doPrintf(warnLevel, warnLevelStr, format, a...)
}

//错误
func (logger *Logger) Error(format string, a ...interface{}) {
	logger.doPrintf(errorLevel, errorLevelStr, format, a...)
}

//致命错误
func (logger *Logger) Fatal(format string, a ...interface{}) {
	logger.doPrintf(fatalLevel, fatalLevelStr, format, a...)
}

var DefaultLogger, _ = New("debug", "", "", log.LstdFlags)

// It's dangerous to call the method on logging
func Export(logger *Logger) {
	if logger != nil {
		DefaultLogger = logger
	}
}

//调试
func Debug(format string, a ...interface{}) {
	DefaultLogger.doPrintf(debugLevel, debugLevelStr, format, a...)
}

//打印信息
func Info(format string, a ...interface{}) {
	DefaultLogger.doPrintf(infoLevel, infoLevelStr, format, a...)
}

//警告
func Warn(format string, a ...interface{}) {
	DefaultLogger.doPrintf(warnLevel, warnLevelStr, format, a...)
}

//错误
func Error(format string, a ...interface{}) {
	DefaultLogger.doPrintf(errorLevel, errorLevelStr, format, a...)
}

//致命错误
func Fatal(format string, a ...interface{}) {
	DefaultLogger.doPrintf(fatalLevel, fatalLevelStr, format, a...)
}

func Close() {
	DefaultLogger.Close()
}
