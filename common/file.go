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

func CreateFile(path string) (fp *os.File, err error) {
	// create dirs if file not exists
	if dir := filepath.Dir(path); dir != "." {
		err = os.MkdirAll(dir, DIR_MODE)
	}
	if err == nil {
		flag := os.O_RDWR|os.O_CREATE|os.O_TRUNC
		fp, err = os.OpenFile(path, flag, FILE_MODE)
	}
	return
}

func OpenFile(path string, readonly, append bool) (fp *os.File, size int64, err error) {
	var exists bool
	size, exists = FileSize(path)
	if size < 0 {
		err = fmt.Errorf("Path is directory or illegal")
		return
	}
	if exists {
		flag := os.O_RDWR
		if readonly {
			flag = os.O_RDONLY
		} else if append {
			flag |= os.O_APPEND
		}
		fp, err = os.OpenFile(path, flag, FILE_MODE)
	} else if readonly == false {
		fp, err = CreateFile(path)
	}
	return
}

// 按分割方法读取文件全部
func ReadFile(path string, split bufio.SplitFunc) ([]string, error) {
	fp, _, err := OpenFile(path, true, false)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	var result []string
	scanner := bufio.NewScanner(fp)
	scanner.Split(split)
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
