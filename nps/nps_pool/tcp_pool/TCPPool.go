package tcp_pool

import "net"

// TCP连接池
type TCPPool struct {
	PoolTCP net.Conn

	// 创建时间(毫秒)
	CreateTime int64
}
