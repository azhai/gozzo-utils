package filesystem

import (
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