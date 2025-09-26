package udp_bridge

				import (
					"DairoNPS/DebugTimer"
				)

import (
	"DairoNPS/dao/dto"
	"DairoNPS/nps"
	"DairoNPS/util/ChannelStatisticsUtil"
	"DairoNPS/util/SecurityUtil"
	"strconv"
	"sync/atomic"
	"time"
)

type UDPBridge struct {
	ClientId        int                                    //客户端ID
	Channel         *dto.ChannelDto                        //频道信息
	ProxyUDPInfo    *nps.UDPInfo                           //代理UDP终端客户端信息
	ClientUDPInfo   *nps.UDPInfo                           //内网穿透客户端UDP信息
	CreateTime      int64                                  // 创建时间(毫秒)
	LastRWTime      int64                                  // 记录最后通信时间(毫秒)
	channelDataSize *ChannelStatisticsUtil.ChannelDataSize //隧道流量统计
}

/**
 * 发送目标端口信息
 */
func (mine *UDPBridge) SendHeaderToClient() error {
DebugTimer.Add199()

	//将加密类型及目标端口 格式:加密状态|端口  1|80   1|127.0.0.1:80
	//1:加密  0:不加密
	head := strconv.Itoa(mine.Channel.SecurityState) + "|" + mine.Channel.TargetPort
	headData := []byte(head)
	err := mine.ClientUDPInfo.Send(headData, len(headData))
	if err != nil {
DebugTimer.Add200()
		return err
	}
	return nil
}

/**
 * 发送数据到内网穿透客户端
 * @param data 要发送的数据数组
 * @param len 要发送的数据长度
 */
func (mine *UDPBridge) SendToClient(data []byte, length int) error {
DebugTimer.Add201()

	//原子递增
	atomic.AddInt64(&mine.channelDataSize.InData, int64(length))
	if mine.Channel.SecurityState == 1 { //加密数据
DebugTimer.Add202()
		SecurityUtil.Mapping(data, length)
	}
	err := mine.ClientUDPInfo.Send(data, length)
	if err != nil {
DebugTimer.Add203()
		return err
	}

	//记录最后通信时间
	mine.LastRWTime = time.Now().UnixMilli()
	return nil
}

/**
 * 通过代理发送数据给终端客户
 * @param data 要发送的数据数组
 * @param len 要发送的数据长度
 */
func (mine *UDPBridge) SendToProxy(data []byte, length int) error {
DebugTimer.Add204()
	if mine.Channel.SecurityState == 1 { //加密数据
DebugTimer.Add205()
		SecurityUtil.Mapping(data, length)
	}
	err := mine.ProxyUDPInfo.Send(data, length)
	if err != nil {
DebugTimer.Add206()
		return err
	}

	//原子递增
	atomic.AddInt64(&mine.channelDataSize.OutData, int64(length))

	//记录最后通信时间
	mine.LastRWTime = time.Now().UnixMilli()
	return nil
}
