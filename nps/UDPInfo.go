package nps

				import (
					"DairoNPS/DebugTimer"
				)

import (
	"fmt"
	"net"
)

/**
 * UDP连接信息
 * mSocket是全局的UDP会话,不能关闭
 */
type UDPInfo struct {
	Udp *net.UDPConn

	//这是客户端地址和端口,并不是服务器端的
	CliAddr *net.UDPAddr
}

/**
 * 连接唯一标识
 */
func (mine *UDPInfo) Key() string {
DebugTimer.Add166()
	return mine.CliAddr.String()
}

/**
 * 向客户端回复数据
 */
func (mine *UDPInfo) Send(data []byte, length int) error {
DebugTimer.Add167()
	_, err := mine.Udp.WriteToUDP(data[:length], mine.CliAddr)
	if err != nil {
DebugTimer.Add168()
		fmt.Println(err)
		return err
	}
	return nil
}
