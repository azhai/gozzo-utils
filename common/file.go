package common

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	FILE_MODE = 0666
	DIR_MODE  = 0777
)

// detect if file exists
// -1, false 不合法的路径
// 0, false 路径不存在
// -1, true 存在文件夹
// >=0, true 文件并存在
func FileSize(path string) (int64, bool) {
	if path == "" {
		return -1, false
	}
	info, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return 0, false
	}
	var size = int64(-1)
	if info.IsDir() == false {
		size = info.Size()
	}
	return size, true
}

func CreateFile(path string, mode os.FileMode) (fp *os.File, err error) {
	// create dirs if file not exists
	if dir := filepath.Dir(path); dir != "." {
		err = os.MkdirAll(dir, DIR_MODE)
	}
	if err == nil {
		fp, err = os.Create(path)
		if err == nil && mode != 0666 {
			fp.Chmod(mode)
		}
	}
	return
}

func OpenFile(path string) (fp *os.File, size int64, err error) {
	var exists bool
	size, exists = FileSize(path)
	if size < 0 {
		err = fmt.Errorf("Path is directory or illegal")
		return
	}
	if exists == false {
		fp, err = CreateFile(path, FILE_MODE)
	} else if size == 0 {
		fp, err = os.OpenFile(path, os.O_RDWR, FILE_MODE)
	} else {
		fp, err = os.Open(path) // 打开方式为 os.O_RDONLY
	}
	return
}

// 按行读取文件全部
func ReadFileLines(path string) ([]string, error) {
	fp, _, err := OpenFile(path)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	return ReadLines(fp)
}

func ReadLines(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
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
