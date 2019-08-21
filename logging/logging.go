package logging

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/azhai/gozzo-utils/filesystem"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 日志等级
var LogLevels = map[string]zapcore.Level{
	"debug":     zapcore.DebugLevel,
	"info":      zapcore.InfoLevel,
	"notice":    zapcore.InfoLevel, // zap中无notice等级
	"warn":      zapcore.WarnLevel,
	"warning":   zapcore.WarnLevel,  // warn的别名
	"err":       zapcore.ErrorLevel, // error的别名
	"error":     zapcore.ErrorLevel,
	"dpanic":    zapcore.DPanicLevel, // Develop环境会panic
	"panic":     zapcore.PanicLevel,  // 都会panic
	"fatal":     zapcore.FatalLevel,
	"critical":  zapcore.FatalLevel, // zap中无critical等级
	"emergency": zapcore.FatalLevel, // zap中无emergency等级
}

var DefaultConfig = LogConfig{
	Encoding:   "console",
	LevelCase:  "capital",
	TimeFormat: "2006-01-02 15:04:05",
	MinLevel:   "info",
	OutputMap:  map[string][]string{":": {"stderr"}},
}

type ILogger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	DPanic(args ...interface{})
	Panic(args ...interface{})
	Fatal(args ...interface{})
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	DPanicf(template string, args ...interface{})
	Panicf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
}

type Logger = zap.SugaredLogger

func NewLogger(level, logdir string) *Logger {
	cfg := DefaultConfig
	cfg.MinLevel = level
	if logdir = strings.TrimSpace(logdir); logdir != "" {
		cfg.OutputMap = map[string][]string{
			":warn":  {filepath.Join(logdir, "access.log")},
			"error:": {filepath.Join(logdir, "error.log")},
		}
	}
	return cfg.BuildSugar()
}

func NewFileLogger(cfg LogConfig, outs map[string]string) *Logger {
	cfg.OutputMap = map[string][]string{}
	for lvl, file := range outs {
		key := fmt.Sprintf("%s:%s", lvl, lvl)
		cfg.OutputMap[key] = []string{file}
	}
	return cfg.BuildSugar()
}

type LogConfig struct {
	Development bool
	Encoding    string
	LevelCase   string
	TimeFormat  string
	MinLevel    string
	OutputMap   map[string][]string
}

func (c LogConfig) BuildSugar() *Logger {
	var cores []zapcore.Core
	encoder := c.BuildEncoder()
	for lvl, outs := range c.OutputMap {
		writer := GetWriteSyncer(outs...)
		pieces := strings.SplitN(lvl, ":", 2)
		stop := strings.TrimSpace(pieces[1])
		start := strings.TrimSpace(pieces[0])
		if start == "" && c.MinLevel != "" {
			start = c.MinLevel
		}
		priority := GetLevelEnabler(start, stop)
		cores = append(cores, zapcore.NewCore(encoder, writer, priority))
	}
	var opts []zap.Option
	if c.Development {
		opts = []zap.Option{zap.Development()}
	}
	logger := zap.New(zapcore.NewTee(cores...), opts...)
	defer logger.Sync()
	return logger.Sugar()
}

func (c LogConfig) BuildEncoder() zapcore.Encoder {
	var encoder zapcore.Encoder
	config := CustomEncoderConfig(c)
	if strings.ToLower(c.Encoding) == "json" {
		encoder = zapcore.NewJSONEncoder(config)
	} else {
		encoder = zapcore.NewConsoleEncoder(config)
	}
	return encoder
}

func CustomEncoderConfig(c LogConfig) zapcore.EncoderConfig {
	ecfg := zap.NewProductionEncoderConfig()
	ecfg.EncodeCaller = nil
	ecfg.EncodeTime = nil
	if c.TimeFormat != "" {
		ecfg.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(c.TimeFormat))
		}
	}
	switch strings.ToLower(c.LevelCase) {
	default:
		ecfg.EncodeLevel = nil
	case "cap", "capital":
		ecfg.EncodeLevel = zapcore.CapitalLevelEncoder
	case "color", "capcolor", "capitalcolor":
		ecfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	case "low", "lower", "lowercase":
		ecfg.EncodeLevel = zapcore.LowercaseLevelEncoder
	case "lowcolor", "lowercolor", "lowercasecolor":
		ecfg.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	}
	return ecfg
}

// 使用绝对路径
func GetLogPath(path string, createIt bool) string {
	absPath, _ := filepath.Abs(path)
	// 解决windows下zap不能识别路径中的盘符问题
	if "windows" == runtime.GOOS {
		re := regexp.MustCompile(`^[A-Za-z]:`)
		absPath = re.ReplaceAllLiteralString(absPath, "")
	}
	if createIt { // 如果不存在就创建文件
		fp, _, err := filesystem.OpenFile(absPath, false, true)
		if err == nil {
			fp.Close()
		}
	}
	return absPath
}

// 日志接收者
func GetWriteSyncer(pathes ...string) zapcore.WriteSyncer {
	if len(pathes) == 0 {
		pathes = []string{""}
	}
	switch pathes[0] {
	case "", "/dev/null":
		return zapcore.AddSync(ioutil.Discard)
	case "stderr":
		return zapcore.Lock(os.Stderr)
	case "stdout":
		return zapcore.Lock(os.Stdout)
	}
	for i, path := range pathes {
		pathes[i] = GetLogPath(path, true)
	}
	sink, closer, err := zap.Open(pathes...)
	if err != nil {
		closer()
		return nil
	}
	return sink
}

// 级别过滤
func GetLevelEnabler(start, stop string) zapcore.LevelEnabler {
	startLvl, stopLvl := zapcore.DebugLevel, zapcore.FatalLevel
	if level, ok := LogLevels[strings.ToLower(start)]; ok {
		startLvl = level
	}
	if level, ok := LogLevels[strings.ToLower(stop)]; ok {
		stopLvl = level
	}
	if stopLvl == zapcore.FatalLevel {
		return zap.NewAtomicLevelAt(startLvl)
	} else if stopLvl == startLvl {
		return zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl == startLvl
		})
	} else {
		return zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= startLvl && lvl <= stopLvl
		})
	}
}
