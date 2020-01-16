package filesystem

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// 当前执行文件所在的目录
func GetRunDir() string {
	return filepath.Dir(os.Args[0])
}

// 取得文件的绝对路径
func GetAbsFile(fname string) string {
	if filepath.IsAbs(fname) == false {
		// 相对于程序运行目录
		dir, err := filepath.Abs(GetRunDir())
		if err != nil {
			return ""
		}
		dir = strings.Replace(dir, "\\", "/", -1)
		fname = filepath.Join(dir, fname)
	}
	return fname
}

// 为文件路径创建目录
func MkdirForFile(path string) int64 {
	size, exists := FileSize(path)
	if size < 0 {
		return size
	}
	if !exists {
		dir := filepath.Dir(path)
		_ = os.MkdirAll(dir, DIR_MODE)
	}
	return size
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

// 遍历目录下的文件
func FindFiles(dir, ext string) (map[string]os.FileInfo, error) {
	var result = make(map[string]os.FileInfo)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return result, err
	}
	for _, file := range files {
		fname := file.Name()
		if ext != "" && !strings.HasSuffix(fname, ext) {
			continue
		}
		fname = filepath.Join(dir, fname)
		result[fname] = file
	}
	return result, nil
}
