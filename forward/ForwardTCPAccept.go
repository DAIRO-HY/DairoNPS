package forward

				import (
					"DairoNPS/DebugTimer"
				)

import (
	"DairoNPS/dao/dto"
	"DairoNPS/util/LogUtil"
	"fmt"
	"net"
	"strings"
)

type ForwardTCPAccept struct {

	//代理端口监听服务
	listen net.Listener

	//转发端口Dto
	forwardDto *dto.ForwardDto

	// 标记监听已经结束
	isFinished bool
}

/**
 * 等待客户端连接
 */
func (mine *ForwardTCPAccept) accept() {
DebugTimer.Add120()
	for {
DebugTimer.Add121()

		//代理服务端Socket
		proxyTCP, err := mine.listen.Accept()
		if err != nil {
DebugTimer.Add122()
			LogUtil.Info(fmt.Sprintf("转发端口:%d 监听结束\n", mine.forwardDto.Port))
			break
		}
		targetIpAndPort := mine.forwardDto.TargetPort
		if !strings.Contains(targetIpAndPort, ":") {
DebugTimer.Add123()
			targetIpAndPort = "127.0.0.1:" + targetIpAndPort
		}

		//目标服务器Socket连接
		targetTCP, err := net.Dial("tcp", targetIpAndPort)
		if err != nil {
DebugTimer.Add124()
			proxyTCP.Close()
			LogUtil.Debug(fmt.Sprintf("转发端口:%d 连接失败\n", mine.forwardDto.Port))
			continue
		}

		//开始桥接
		startBridge(mine.forwardDto, proxyTCP, targetTCP)
	}
	mine.listen.Close()
	LogUtil.Debug(fmt.Sprintf("转发端口:%d 监听结束\n", mine.forwardDto.Port))
	mine.isFinished = true
}

/**
 * 停止监听端口
 */
func (mine *ForwardTCPAccept) shutdown() {
DebugTimer.Add125()
	mine.listen.Close()

	//关闭当前的桥接通信
	shutdownBridge(mine.forwardDto.Id)
	removeAccept(mine.forwardDto.Id)
}
