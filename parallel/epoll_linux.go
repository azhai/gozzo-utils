// +build linux

package parallel

import (
	"syscall"
)

type Poll struct {
	fd    int // epoll fd
	wfd   int // wake fd
}

func NewPoll() *Poll {
	fd, err := syscall.EpollCreate1(0)
	if err != nil {
		panic(err)
	}
	r1, _, eno := syscall.Syscall(syscall.SYS_EVENTFD2, 0, 0, 0)
	if eno != 0 {
		syscall.Close(fd)
		panic(err)
	}
	p := &Poll{fd: fd, wfd: int(r1)}
	p.AddRead(p.wfd)
	return p
}

func (p *Poll) Close() error {
	if err := syscall.Close(p.wfd); err != nil {
		return err
	}
	return syscall.Close(p.fd)
}

func (p *Poll) Wait(action func(fd int, flags uint32) error) error {
	var fd int
	events := make([]syscall.EpollEvent, 64)
	for {
		n, err := syscall.EpollWait(p.fd, events, -1)
		if err != nil && err != syscall.EINTR {
			return err
		}
		for i := 0; i < n; i++ {
			if fd = int(events[i].Fd); fd == p.wfd {
				continue
			}
			err := action(fd, events[i].Events)
			if err != nil {
				return err
			}
		}
	}
}

func (p *Poll) Trigger(note interface{}) error {
	_, err := syscall.Write(p.wfd, []byte{0, 0, 0, 0, 0, 0, 0, 1})
	return err
}

func (p *Poll) Control(fd, op int, flags uint32) error {
	ev := &syscall.EpollEvent{Fd: int32(fd), Events: flags}
	return syscall.EpollCtl(p.fd, op, fd, ev)
}

func (p *Poll) AddRead(fd int) error {
	return p.Control(fd, syscall.EPOLL_CTL_ADD, syscall.EPOLLIN)
}

func (p *Poll) DelRead(fd int) error {
	return p.Control(fd, syscall.EPOLL_CTL_DEL, syscall.EPOLLIN)
}

func (p *Poll) AddReadWrite(fd int) error {
	return p.Control(fd, syscall.EPOLL_CTL_ADD, syscall.EPOLLIN|syscall.EPOLLOUT)
}

func (p *Poll) DelReadWrite(fd int) error {
	return p.Control(fd, syscall.EPOLL_CTL_DEL, syscall.EPOLLIN|syscall.EPOLLOUT)
}
