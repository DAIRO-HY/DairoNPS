package forward

				import (
					"DairoNPS/DebugTimer"
				)

import (
	"DairoNPS/dao/ForwardDao"
	"DairoNPS/dao/dto"
	"DairoNPS/util/ForwardStatisticsUtil"
	"DairoNPS/util/LogUtil"
	"fmt"
	"net"
	"strconv"
	"sync"
)

/**
 * 数据转发连接管理
 */

/**
 * 数据转发id对应的端口监听
 */
var forwardAcceptMap = make(map[int]*ForwardTCPAccept)

/**
 * forwardAcceptMap操作互斥锁
 */
var forwardAcceptLock sync.Mutex

// 端口转发代理端口数量
func GetAcceptCount() int {
DebugTimer.Add126()
	count := 0
	forwardAcceptLock.Lock()
	count = len(forwardAcceptMap)
	forwardAcceptLock.Unlock()
	return count
}

/**
 * 开启端口监听
 */
func StartAcceptAll() {
DebugTimer.Add127()

	//开启NPS客户端ID下所有的隧道
	list := ForwardDao.SelectActive()
	for _, item := range list {
DebugTimer.Add128()
		Accept(item)
	}
}

/**
 * 开启端口监听
 * @param forwardDto 隧道信息
 */
func Accept(forwardDto *dto.ForwardDto) {
DebugTimer.Add129()

	//加载统计数据
	ForwardStatisticsUtil.Init()
	forwardAcceptLock.Lock()
	old := forwardAcceptMap[forwardDto.Id]
	if old != nil { //若该隧道已经在监听,则先停止
DebugTimer.Add130()
		old.shutdown()
	}

	listen, err := net.Listen("tcp", ":"+strconv.Itoa(forwardDto.Port))
	if err != nil {
DebugTimer.Add131()
		errMsg := fmt.Sprintf("转发端口:%d 监听失败。err:%q\n", forwardDto.Port, err)
		ForwardDao.SetError(forwardDto.Id, &errMsg)
		LogUtil.Debug(errMsg)
		forwardAcceptLock.Unlock()
		return
	}
	ForwardDao.SetError(forwardDto.Id, nil)
	LogUtil.Info(fmt.Sprintf("转发端口:%d 监听开始", forwardDto.Port))
	tcpAccept := &ForwardTCPAccept{
		forwardDto: forwardDto,
		listen:     listen,
	}

	forwardAcceptMap[forwardDto.Id] = tcpAccept

	//开启监听
	go tcpAccept.accept()
	forwardAcceptLock.Unlock()
}

/**
 * 移除隧道监听列表
 */
func removeAccept(forwardId int) {
DebugTimer.Add132()
	forwardAcceptLock.Lock()
	tcpAccept := forwardAcceptMap[forwardId]
	if tcpAccept == nil {
DebugTimer.Add133()
		forwardAcceptLock.Unlock()
		return
	}
	delete(forwardAcceptMap, forwardId)

	////关闭隧道的时候保存流量
	//ForwardDao.setDataLen(proxy.forwardDto)
	forwardAcceptLock.Unlock()
}

/**
 * 关闭监听
 * @param forwardId 隧道id
 */
func CloseAccept(forwardId int) {
DebugTimer.Add134()
	forwardAcceptLock.Lock()
	tcpAccept := forwardAcceptMap[forwardId]
	forwardAcceptLock.Unlock()
	if tcpAccept != nil {
DebugTimer.Add135()
		tcpAccept.shutdown()
	}
}
