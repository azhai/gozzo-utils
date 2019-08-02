package logging

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/azhai/gozzo-utils/common"
	"go.uber.org/zap"
)

var (
	lang    = "open source"
	soft    = "simple, fast, and reliable"
	tpl     = "Go is an %s programming language, designed for building %s software."
	msg     = fmt.Sprintf(tpl, lang, soft)
	logdir = "/tmp/sugar_logger_test/"
	outfile = "/tmp/sugar_logger_test/access.log"
	errfile = "/tmp/sugar_logger_test/error.log"
)

func CreateLogger(level string) *zap.SugaredLogger {
	os.Truncate(outfile, 0)
	return NewLogger(level, logdir).Named("Test")
}

func ReadLastLog(logfile string) (level string, content string) {
	data, _ := common.ReadFileTail(logfile, 1024)
	data = bytes.TrimRight(data, "\r\n")
	n := bytes.LastIndexByte(data, byte('\n'))
	pieces := bytes.SplitN(data[n+1:], []byte("\t"), 4)
	if len(pieces) >= 4 {
		level, content = string(pieces[1]), string(pieces[3])
	}
	return
}

// 测试不记录的 INFO 级别
func TestInfo(t *testing.T) {
	logger := CreateLogger("warning")
	logger.Infof(tpl, lang, soft)
	level, content := ReadLastLog(outfile)
	t.Log(level)
	t.Log(content)
	if level == "INFO" {
		t.Errorf("The level is %s", level)
	} else if content != "" {
		t.Errorf("The content is %s", content)
	}
}

// 测试要记录的 WARN 级别
func TestWarn(t *testing.T) {
	logger := CreateLogger("warn")
	logger.Warn(msg)
	logger.Warnf(tpl, lang, soft)
	level, content := ReadLastLog(outfile)
	t.Log(level)
	t.Log(content)
	if level != "WARN" {
		t.Errorf("The level is %s", level)
	} else if content != msg {
		t.Errorf("The content is %s", content)
	}
}

// 记录到另一个文件的 ERROR 级别
func TestError(t *testing.T) {
	logger := CreateLogger("warning")
	logger.Errorf(tpl, lang, soft)
	level, content := ReadLastLog(errfile)
	t.Log(level)
	t.Log(content)
	if level == "ERROR" {
		t.Errorf("The level is %s", level)
	} else if content != "" {
		t.Errorf("The content is %s", content)
	}
}

// 不记录
func BenchmarkEmpty(b *testing.B) {
	logger := CreateLogger("warning")
	for i := 0; i < b.N; i++ {
		logger.Infof(tpl, lang, soft)
	}
}

// 记录到文件，但不格式化字符串
func BenchmarkFile1(b *testing.B) {
	logger := CreateLogger("warning")
	for i := 0; i < b.N; i++ {
		logger.Warn(msg)
	}
}

// 记录到文件
func BenchmarkFile2(b *testing.B) {
	logger := CreateLogger("warning")
	tpl = "%d " + msg
	for i := 0; i < b.N; i++ {
		logger.Warnf(tpl, i)
	}
}
