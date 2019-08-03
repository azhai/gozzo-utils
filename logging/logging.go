package logging

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/azhai/gozzo-utils/common"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger = zap.SugaredLogger

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
	LeastLevel: "info",
	TimeFormat: "2006-01-02 15:04:05",
	Outputs:    []string{"stderr"},
}

func NewLogger(level, logdir string) *Logger {
	cfg := DefaultConfig
	cfg.LeastLevel = level
	if logdir = strings.TrimSpace(logdir); logdir != "" {
		cfg.ErrorFile = filepath.Join(logdir, "error.log")
		cfg.Outputs = []string{filepath.Join(logdir, "access.log")}
	}
	return cfg.BuildSugar()
}

type LogConfig struct {
	Development bool
	Encoding    string
	LeastLevel  string
	TimeFormat  string
	ErrorFile   string
	Outputs     []string
}

func (c LogConfig) BuildSugar() *Logger {
	var encoder zapcore.Encoder
	config := CustomEncoderConfig(c.TimeFormat)
	if strings.ToLower(c.Encoding) == "json" {
		encoder = zapcore.NewJSONEncoder(config)
	} else {
		encoder = zapcore.NewConsoleEncoder(config)
	}

	var cores []zapcore.Core
	writer := GetWriteSyncer(c.Outputs...)
	priority := GetLevelEnabler(c.LeastLevel, "")
	if c.ErrorFile != "" {
		errWriter := GetWriteSyncer(c.ErrorFile)
		errPriority := GetLevelEnabler("error", "")
		priority = GetLevelEnabler(c.LeastLevel, "warn")
		cores = append(cores, zapcore.NewCore(encoder, errWriter, errPriority))
	}
	cores = append(cores, zapcore.NewCore(encoder, writer, priority))

	var opts []zap.Option
	if c.Development {
		opts = []zap.Option{zap.Development()}
	}
	logger := zap.New(zapcore.NewTee(cores...), opts...)
	defer logger.Sync()
	return logger.Sugar()
}

func CustomEncoderConfig(layout string) zapcore.EncoderConfig {
	ecfg := zap.NewProductionEncoderConfig()
	ecfg.EncodeCaller = nil
	ecfg.EncodeLevel = zapcore.CapitalLevelEncoder
	ecfg.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(layout))
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
		fp, _, err := common.OpenFile(absPath, false, true)
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
