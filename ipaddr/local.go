package ipaddr

import (
	"net"
	"sort"

	"github.com/azhai/gozzo-utils/common"
)

type LocalIP struct {
	Locals []net.IP
}

// 字符串转为IP
func NewLocalIP(hosts ...string) *LocalIP {
	l := &LocalIP{}
	for _, host := range hosts {
		l.Locals = append(l.Locals, net.ParseIP(host))
	}
	return l
}

func (l *LocalIP) Count() int {
	if len(l.Locals) == 0 {
		l.Locals = GetLocalList()
	}
	return len(l.Locals)
}

// 转为UDP地址
func (l *LocalIP) ToUDP(ip net.IP) *net.UDPAddr {
	return &net.UDPAddr{IP: ip}
}

// 转为TCP地址
func (l *LocalIP) ToTCP(ip net.IP) *net.TCPAddr {
	return &net.TCPAddr{IP: ip}
}

func (l *LocalIP) GetIPs(offset, limit int) []net.IP {
	count := l.Count()
	start, stop := common.GetStartStop(offset, limit, count)
	if start < 0 || stop <= 0 {
		return nil
	}
	return l.Locals[start:stop]
}

// 获取所有局域网IP，除了127.0.0.1
func GetLocalList() []net.IP {
	var locals []net.IP
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return locals
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok {
			if ipnet.IP.IsLoopback() || ipnet.IP.To4() == nil {
				continue
			}
			locals = append(locals, ipnet.IP)
		}
	}
	sort.Slice(locals, func(i, j int) bool {
		return locals[i].String() < locals[j].String()
	})
	return locals
}

// 给出几个本机局域网TCP地址
func GetLocalAddrs(limit, offset int) []*net.TCPAddr {
	if limit <= 1 && offset == 0 {
		return make([]*net.TCPAddr, 1)
	}
	locals := NewLocalIP().GetIPs(offset, limit)
	var laddrs []*net.TCPAddr
	for _, ip := range locals {
		laddrs = append(laddrs, &net.TCPAddr{IP: ip})
	}
	return laddrs
}
