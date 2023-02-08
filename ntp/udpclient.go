package ntp

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
)

// Udp客户端
type UdpClient struct {
	addr        string       // 地址
	conn        *net.UDPConn // 连接
	size        int          // 大小
	ReceiveChan chan *Packet // 接收通道
	SendChan    chan *Packet // 发送通道
}

// 构造函数
func NewUdpClient(host string, port uint) *UdpClient {
	return &UdpClient{
		addr:        fmt.Sprintf("%s:%d", host, port),
		size:        65535 - 20 - 8, // IPv4 max size - IPv4 Header size - UDP Header size
		ReceiveChan: make(chan *Packet, 10),
		SendChan:    make(chan *Packet, 10),
	}
}

func (c *UdpClient) Open() error {
	defer c.Close()

	// 监听
	addr, err := net.ResolveUDPAddr("udp", c.addr)
	if err != nil {
		return err
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return err
	}
	c.conn = conn
	// 写
	go func() {
		for {
			p := <-c.SendChan
			_, err = c.conn.WriteTo(p.Data, p.Addr)
			if err != nil {
				logrus.Error(err)
			}
		}
	}()
	// 读
	for {
		data := make([]byte, c.size)
		n, addr, err := conn.ReadFromUDP(data)
		if err != nil {
			logrus.Error(err)
			continue
		}
		c.ReceiveChan <- &Packet{
			Addr: addr,
			Data: data[:n],
		}
	}
}

func (c *UdpClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}

	return nil
}

func (c *UdpClient) Receive() <-chan *Packet {
	return c.ReceiveChan
}

func (c *UdpClient) Send(packet *Packet) {
	c.SendChan <- packet
}

func (c *UdpClient) Heartbeat(packet *Packet) {
	if packet == nil {
		packet = &Packet{
			Data: []byte("ping"),
		}
	}
	c.Send(packet)
}
