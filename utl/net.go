package utl

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

// 获取Ip地址
func GetIp() (string, error) {
	var ip string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ip, err
	}
	for _, a := range addrs {
		if ipNet, ok := a.(*net.IPNet); ok {
			if !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
				ip = ipNet.IP.String()
				break
			}
		}
	}

	return ip, nil
}

// 获取端口
func GetPort(min, max int) int {
	for i := 0; i < 5; i++ {
		port := min + rand.Intn(max-min)
		_, err := net.DialTimeout("udp", fmt.Sprintf(":%d", port), 3*time.Second)
		if err != nil {
			return port
		}
		continue
	}

	return 0
}
