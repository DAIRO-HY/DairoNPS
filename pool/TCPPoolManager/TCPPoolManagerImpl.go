package TCPPoolManager

import "net"

// 暴露接口给外部调用，防止循环引用
var GetAndAddPool func(clientID int) net.Conn = getAndAddPool
