package tcp_pool

import "net"

// 客户端连接池
type TCPPool struct {
	ClientID int

	PoolTCP net.Conn

	// 创建时间(毫秒)
	CreateTime int64
}
