package udp_bridge

import (
	"DairoNPS/dao/dto"
	"DairoNPS/nps"
	"DairoNPS/util/SecurityUtil"
	"strconv"
)

/**
 * UDP桥接管理
 * @param client 客户端DTO
 * @param channel 频道信息
 * @param proxyInfo 代理UDP终端客户端信息
 * @param clientInfo 内网穿透客户端UDP信息
 */
type UDPBridge struct {
	client  *dto.ClientDto
	channel *dto.ChannelDto

	proxyUDPInfo  *nps.UDPInfo
	clientUDPInfo *nps.UDPInfo

	/**
	 * 最后一次读取到数据的时间,用来判断Socket是否存活
	 */
	//var lastSessionTime = System.currentTimeMillis()

	/**
	 * 是否加密数据
	 */
	isEncodeData bool
}

func (mine *UDPBridge) start() {
	mine.sendHeaderToClient()
}

/**
 * 发送目标端口信息
 */
func (mine *UDPBridge) sendHeaderToClient() {

	//将加密类型及目标端口 格式:加密状态|端口  1|80   1|127.0.0.1:80
	//1:加密  0:不加密
	head := strconv.Itoa(mine.channel.SecurityState) + "|" + mine.channel.TargetPort
	headData := []byte(head)
	err := mine.clientUDPInfo.Send(headData, len(headData))
	if err != nil {
		return
	}
}

/**
 * 发送数据到内网穿透客户端
 * @param data 要发送的数据数组
 * @param len 要发送的数据长度
 */
func (mine *UDPBridge) SendToClient(data []byte, length int) {

	//记录最后通信时间
	//this.lastSessionTime = System.currentTimeMillis()
	if mine.isEncodeData { //加密数据
		SecurityUtil.Mapping(data, length)
	}
	err := mine.clientUDPInfo.Send(data, length)
	if err != nil {
		return
	}

	//入网统计
	//this.inDataTotal = this.inDataTotal.plus(len)
	//this.channel.inDataTotal = this.channel.inDataTotal!!.plus(len)
	//this.client.inDataTotal = this.client.inDataTotal!!.plus(len)
	//CLSConfig.systemConfig.inDataTotal = CLSConfig.systemConfig.inDataTotal!!.plus(len)
}

/**
 * 通过代理发送数据给终端客户
 * @param data 要发送的数据数组
 * @param len 要发送的数据长度
 */
func (mine *UDPBridge) SendToProxy(data []byte, length int) {
	if mine.isEncodeData { //加密数据
		SecurityUtil.Mapping(data, length)
	}
	err := mine.proxyUDPInfo.Send(data, length)
	if err != nil {
		return
	}
}
