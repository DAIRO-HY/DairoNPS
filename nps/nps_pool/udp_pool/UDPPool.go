package udp_pool

				import (
					"DairoNPS/DebugTimer"
				)

import (
	"DairoNPS/constant/NPSConstant"
	"DairoNPS/nps"
)

// UDP连接池
type UDPPool struct {
	UDPInfo *nps.UDPInfo

	// 创建时间(毫秒)
	CreateTime int64
}

// 通知客户端关闭该连接池
func (mine *UDPPool) CloseNotify() {
DebugTimer.Add314()
	closeData := []byte(NPSConstant.UDP_POOL_CLOSE_FLAG)
	mine.UDPInfo.Send(closeData, len(closeData))
}
