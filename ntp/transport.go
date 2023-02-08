package ntp

import (
	"net"
)

// 包
type Packet struct {
	Addr net.Addr // 地址
	Data []byte   // 数据
}

// 传输接口
type ITransport interface {
	Open() error             // 打开
	Close() error            // 关闭
	Receive() <-chan *Packet // 接收
	Send(packet *Packet)     // 发送
}

// 服务器接口
type IServer interface {
	ITransport
}

// 客户端接口
type IClient interface {
	ITransport
	Heartbeat(packet *Packet) // 心跳
}

// 构造函数
func NewServer(network string, port uint16) IServer {
	if network == "udp" {
		return NewUdpServer(port)
	}

	return nil
}
