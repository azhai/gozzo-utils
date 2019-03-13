package logging

import (
	"path/filepath"
	"strings"
	"time"

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

func GetLogPath(path string) string {
	if path == "" || path == "/dev/null" {
		return "/dev/null"
	}
	if path == "stdout" || path == "stderr" {
		return path
	}
	absPath, _ := filepath.Abs(path)
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

func NewLogger(level, errput string, outputs ...string) *zap.SugaredLogger {
	cfg := NewConfig(errput, outputs...)
	cfg.SetLevelName(level)
	return cfg.BuildSugar()
}

func NewDevLogger(errput string, outputs ...string) *zap.SugaredLogger {
	cfg := NewConfig(errput, outputs...)
	cfg.SetDevMode(true)
	cfg.SetLevelName("debug")
	cfg.SetTimeFormat("2006-01-02 15:04:05.999")
	return cfg.BuildSugar()
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

func (c *LogConfig) BuildSugar() *zap.SugaredLogger {
	if logger, err := c.Build(); err == nil {
		return logger.Sugar()
	}
	return nil
}
