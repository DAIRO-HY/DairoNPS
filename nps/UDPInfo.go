package nps

import "net"

/**
 * UDP连接信息
 * mSocket是全局的UDP会话,不能关闭
 */
type UDPInfo struct {
	Socket *net.UDPConn

	//这是客户端地址和端口,并不是服务器端的
	Addr *net.UDPAddr
}

/**
 * 连接唯一标识
 */
func (mine *UDPInfo) Key() string {
	//"${this.inet.hostAddress}:$port"
	return mine.Addr.String()
}

/**
 * 向客户端回复数据
 */
func (mine *UDPInfo) Send(data []byte, length int) error {
	_, err := mine.Socket.WriteToUDP(data[:length], mine.Addr)
	if err != nil {
		mine.Socket.Close()
		return err
	}
	return nil
}
