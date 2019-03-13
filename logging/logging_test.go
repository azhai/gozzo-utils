package logging

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"go.uber.org/zap"
)

var (
	lang    = "open source"
	soft    = "simple, fast, and reliable"
	tpl     = "Go is an %s programming language, designed for building %s software."
	msg     = fmt.Sprintf(tpl, lang, soft)
	logfile = "/tmp/sugar_logger_test.log"
)

func CreateLogger(level string) *zap.SugaredLogger {
	os.Truncate(logfile, 0)
	return NewLogger(level, logfile, logfile).Named("Test")
}

// 读取文件末尾若干字节
func ReadFileTail(file string, size int) (data []byte, err error) {
	var fp *os.File
	if fp, err = os.Open(file); err != nil {
		return
	}
	defer fp.Close()
	offset, err := fp.Seek(0-int64(size), os.SEEK_END)
	// 当size超出文件大小时，游标移到开头并报错，这里忽略错误
	if offset >= 0 {
		data = make([]byte, size)
		reads, err := fp.Read(data)
		if reads >= 0 {
			data = data[:reads]
		}
		if err == io.EOF {
			err = nil
		}
	}
	return
}

func ReadLastLog(t *testing.T) (level string, content string) {
	data, _ := ReadFileTail(logfile, 1024)
	data = bytes.TrimRight(data, "\r\n")
	n := bytes.LastIndexByte(data, byte('\n'))
	pieces := bytes.SplitN(data[n+1:], []byte("\t"), 4)
	if len(pieces) >= 4 {
		level, content = string(pieces[1]), string(pieces[3])
	}
	t.Log(level)
	t.Log(content)
	return
}

// 测试要记录的 WARN 级别
func TestWarn(t *testing.T) {
	logger := CreateLogger("warn")
	logger.Warn(msg)
	logger.Warnf(tpl, lang, soft)
	level, content := ReadLastLog(t)
	if level != "WARN" {
		t.Errorf("The level is %s", level)
	} else if content != msg {
		t.Errorf("The content is %s", content)
	}
}

// 测试不记录的 INFO 级别
func TestInfo(t *testing.T) {
	logger := CreateLogger("warning")
	logger.Infof(tpl, lang, soft)
	level, content := ReadLastLog(t)
	if level == "INFO" {
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
