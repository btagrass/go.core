package ntp

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
)

// Udp服务器
type UdpServer struct {
	addr        string       // 地址
	conn        *net.UDPConn // 连接
	size        int          // 大小
	ReceiveChan chan *Packet // 接收通道
	SendChan    chan *Packet // 发送通道
}

// 构造函数
func NewUdpServer(port uint16) *UdpServer {
	return &UdpServer{
		addr:        fmt.Sprintf(":%d", port),
		size:        65535 - 20 - 8, // IPv4 max size - IPv4 Header size - UDP Header size
		ReceiveChan: make(chan *Packet, 10),
		SendChan:    make(chan *Packet, 10),
	}
}

func (s *UdpServer) Open() error {
	defer s.Close()

	// 监听
	addr, err := net.ResolveUDPAddr("udp", s.addr)
	if err != nil {
		return err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	s.conn = conn
	// 写
	go func() {
		for {
			select {
			case p := <-s.SendChan:
				_, err = s.conn.WriteTo(p.Data, p.Addr)
				if err != nil {
					logrus.Error(err)
				}
			}
		}
	}()
	// 读
	for {
		data := make([]byte, s.size)
		n, addr, err := conn.ReadFromUDP(data)
		if err != nil {
			logrus.Error(err)
			continue
		}
		s.ReceiveChan <- &Packet{
			Addr: addr,
			Data: data[:n],
		}
	}
}

func (s *UdpServer) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}

	return nil
}

func (s *UdpServer) Receive() <-chan *Packet {
	return s.ReceiveChan
}

func (s *UdpServer) Send(packet *Packet) {
	s.SendChan <- packet
}
