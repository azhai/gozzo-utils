package logging

import (
	"github.com/azhai/gozzo-utils/common"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

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

func GetLogPath(path string) string {
	if path == "" || path == "/dev/null" {
		return "/dev/null"
	}
	if path == "stdout" || path == "stderr" {
		return path
	}
	absPath, _ := filepath.Abs(path)
	// 解决windows下zap不能识别路径中的盘符问题
	if "windows" == runtime.GOOS {
		re := regexp.MustCompile(`^[A-Za-z]:`)
		absPath = re.ReplaceAllLiteralString(absPath, "")
	}
	return absPath
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

func NewLogger(level, errput string, outputs ...string) *Logger {
	cfg := NewConfig(errput, outputs...)
	cfg.SetLevelName(level)
	return cfg.BuildSugar()
}

func NewDevLogger(errput string, outputs ...string) *Logger {
	cfg := NewConfig(errput, outputs...)
	cfg.SetDevMode(true)
	cfg.SetLevelName("debug")
	cfg.SetTimeFormat("2006-01-02 15:04:05.999")
	return cfg.BuildSugar()
}

func NewLoggerInDir(level, logdir string) *Logger {
	errfile := filepath.Join(logdir, "error.log")
	if fp, _, err := common.OpenFile(errfile, false, true); err == nil {
		fp.Close()
	}
	outfile := filepath.Join(logdir, "access.log")
	if fp, _, err := common.OpenFile(outfile, false, true); err == nil {
		fp.Close()
	}
	return NewLogger(level, errfile, outfile)
}

type LogConfig struct {
	zap.Config
	LevelName  string
	TimeFormat string
}

func NewConfig(errput string, outputs ...string) *LogConfig {
	cfg := zap.NewProductionConfig()
	cfg.ErrorOutputPaths = []string{GetLogPath(errput)}
	cfg.OutputPaths = make([]string, 0) // 清空，原来是 []string{"stderr"}
	for _, output := range outputs {
		cfg.OutputPaths = append(cfg.OutputPaths, GetLogPath(output))
	}
	cfg.Encoding = "console"
	cfg.EncoderConfig = CustomEncoderConfig("2006-01-02 15:04:05")
	return &LogConfig{cfg, "info", "2006-01-02 15:04:05"}
}

func (c *LogConfig) SetDevMode(isDev bool) {
	c.Development = isDev
}

func (c *LogConfig) SetLevelName(level string) {
	c.LevelName = strings.ToLower(level)
	if lvl, ok := LogLevels[c.LevelName]; ok { // 默认INFO及以上
		c.Level = zap.NewAtomicLevelAt(lvl)
	}
}

func (c *LogConfig) SetTimeFormat(layout string) {
	c.TimeFormat = layout
	c.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(c.TimeFormat))
	}
}

func (c *LogConfig) BuildSugar() *Logger {
	logger, err := c.Build()
	if err != nil {
		panic(err)
		return nil
	}
	return logger.Sugar()
}
