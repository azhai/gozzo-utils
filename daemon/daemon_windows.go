// +build windows

package daemon

import "github.com/kardianos/service"

/*
如果程序要作为Windows服务启动，需要将原来的main()改名为run()
新加一个文件main_windows.go，开头加入和本文件一样的build tag，导入本包
然后参考下面的程序写法：
func main() {
	name := "My_Example_Service"
	desc := "一个样例服务（重要勿删）"
	daemon.WinMain(name, desc, run)
}

添加Windows安装服务脚本如下：（注意Windows系统有bug，所以=后面要留有空格）
sc create "My_Example_Service" displayName= "My_Example_Service" binPath= "D:\gocode\my_service.exe" start= auto
sc description "My_Example_Service" "一个样例服务（重要勿删）"
PAUSE
*/

func WinMain(name, desc string, run func()) {
	info := NewConfig(name, desc)
	prg := &Program{Main: run}
	s, _ := service.New(prg, info)
	s.Run()
}

func NewConfig(name, desc string) *service.Config {
	return &service.Config{
		Name:        name,
		DisplayName: name,
		Description: desc,
	}
}

type Program struct {
	Main func()
}

func (p *Program) Start(s service.Service) error {
	go p.Main()
	return nil
}

func (p *Program) Stop(s service.Service) error {
	return nil
}
