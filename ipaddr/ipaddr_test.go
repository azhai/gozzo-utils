package ipaddr

import (
	"fmt"
	"net"
	"testing"
)

func GetLocalIP() *LocalIP {
	ip := NewLocalIP()
	for i := 70; i < 75; i++ {
		loc := net.ParseIP(fmt.Sprintf("192.168.0.%d", i))
		ip.Locals = append(ip.Locals, loc)
	}
	return ip
}

func TestListLocal(t *testing.T) {
	t.Log(GetLocalList())
}

func TestListSlice(t *testing.T) {
	ip := GetLocalIP()
	t.Log(0, 0, ip.GetIPs(0, 0))
	t.Log(0, 1, ip.GetIPs(0, 1))
	t.Log(0, -1, ip.GetIPs(0, -1))
	t.Log(1, 0, ip.GetIPs(1, 0))
	t.Log(1, 1, ip.GetIPs(1, 1))
	t.Log(1, -1, ip.GetIPs(1, -1))
	t.Log(-1, 0, ip.GetIPs(-1, 0))
	t.Log(-1, 1, ip.GetIPs(-1, 1))
	t.Log(-1, -1, ip.GetIPs(-1, -1))
}
