package filesystem

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
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
		flag := os.O_RDWR | os.O_CREATE | os.O_TRUNC
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

// CopyFile copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

// 通过Bash命令复制整个目录，只能运行于Linux或MacOS
// 当dst结尾带斜杠时，复制为dst下的子目录
func CopyDir(src, dst string) (err error) {
	if length := len(src); src[length-1] == '/' {
		src = src[:length-1] //去掉结尾的斜杠
	}
	info, err := os.Stat(src)
	if err != nil || !info.IsDir() {
		return
	}
	err = exec.Command("cp", "-rf", src, dst).Run()
	return
}

// 使用 wc -l 计算有多少行
func LineCount(fname string) int {
	var num int
	fname, err := filepath.Abs(fname)
	if err != nil {
		return -1
	}
	out, err := exec.Command("wc", "-l", fname).Output()
	if err != nil {
		return -1
	}
	col := strings.SplitN(string(out), " ", 2)[0]
	if num, err = strconv.Atoi(col); err != nil {
		return -1
	}
	return num
}
