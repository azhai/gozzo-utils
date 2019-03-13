// +build windows

package daemon

import "github.com/kardianos/service"

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

func (this *Program) Start(s service.Service) error {
	go this.Main()
	return nil
}

func (this *Program) Stop(s service.Service) error {
	return nil
}
