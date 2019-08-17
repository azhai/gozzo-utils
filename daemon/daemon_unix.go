// +build android darwin dragonfly freebsd linux netbsd openbsd solaris

package daemon

/*
本文件用作说明unix系统的build tag
请先将原来的main()改名为run()，新建main_unix.go文件并加入如下代码：

func main() {
	run()
}
*/
