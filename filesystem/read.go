package filesystem

import (
	"bufio"
	"io"
)

// 每次只保留当前行数据
type LineReader struct {
	err  error
	line []byte
	rd   io.ReadCloser
	*bufio.Reader
}

func NewLineReader(fname string) *LineReader {
	fp, _, err := OpenFile(fname, true, false)
	return &LineReader{err: err, rd: fp, Reader: bufio.NewReader(fp)}
}

func (r *LineReader) Close() error {
	r.err = r.rd.Close()
	return r.err
}

func (r *LineReader) Err() error {
	return r.err
}

func (r *LineReader) Line() []byte {
	return r.line
}

func (r *LineReader) Text() string {
	return string(r.line)
}

func (r *LineReader) Reading() bool {
	line, isPrefix, err := r.ReadLine()
	if isPrefix == false {
		r.line = line
	} else if line != nil {
		r.line = append(r.line, line...)
	}
	if err == io.EOF {
		return false
	} else if err != nil {
		r.err = err
	}
	return true
}

// 读取全部数据，按行组成列表
func ReadLines(path string) ([]string, error) {
	fp, _, err := OpenFile(path, true, false)
	defer fp.Close()
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(fp)
	scanner.Split(bufio.ScanLines)
	var result []string
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}
	return result, scanner.Err()
}

// 读取文件末尾若干字节
func ReadFileTail(path string, size int) ([]byte, error) {
	fp, _, err := OpenFile(path, true, false)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	offset, err := fp.Seek(0-int64(size), io.SeekEnd)
	if offset < 0 {
		return nil, err
	}
	// 当size超出文件大小时，游标移到开头并报错，这里忽略错误
	result := make([]byte, size)
	reads, err := fp.Read(result)
	if reads >= 0 {
		result = result[:reads]
	}
	if err == io.EOF {
		err = nil
	}
	return result, err
}
