package forward

import (
	"DairoNPS/dao/ForwardDao"
	"DairoNPS/dao/dto"
	"DairoNPS/util/ForwardStatisticsUtil"
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
var forwardIdToForwardAccept = make(map[int]*ForwardTCPAccept)

/**
 * forwardIdToForwardAccept操作互斥锁
 */
var forwardIdToForwardAcceptLock sync.Mutex

/**
 * 开启端口监听
 */
func StartAcceptAll() {

	//开启NPS客户端ID下所有的隧道
	list := ForwardDao.SelectActive()
	for _, item := range list {
		Accept(item)
	}
}

/**
 * 开启端口监听
 * @param forwardDto 隧道信息
 */
func Accept(forwardDto *dto.ForwardDto) {

	//加载统计数据
	ForwardStatisticsUtil.Init()
	forwardIdToForwardAcceptLock.Lock()
	old := forwardIdToForwardAccept[forwardDto.Id]
	if old != nil { //若该隧道已经在监听,则先停止
		old.shutdown()
	}

	listen, err := net.Listen("tcp", ":"+strconv.Itoa(forwardDto.Port))
	if err != nil {
		fmt.Printf("转发端口:%d 监听失败。err:%p\n", forwardDto.Port, err)
		forwardIdToForwardAcceptLock.Unlock()
		return
	}
	fmt.Printf("转发端口:%d 监听开始", forwardDto.Port)
	tcpAccept := &ForwardTCPAccept{
		forwardDto: forwardDto,
		listen:     listen,
	}

	forwardIdToForwardAccept[forwardDto.Id] = tcpAccept

	//开启监听
	go tcpAccept.accept()
	forwardIdToForwardAcceptLock.Unlock()
}

/**
 * 移除隧道监听列表
 */
func removeAccept(forwardId int) {
	forwardIdToForwardAcceptLock.Lock()
	tcpAccept := forwardIdToForwardAccept[forwardId]
	if tcpAccept == nil {
		forwardIdToForwardAcceptLock.Unlock()
		return
	}
	delete(forwardIdToForwardAccept, forwardId)

	////关闭隧道的时候保存流量
	//ForwardDao.setDataLen(proxy.forwardDto)
	forwardIdToForwardAcceptLock.Unlock()
}

/**
 * 关闭监听
 * @param forwardId 隧道id
 */
func CloseAccept(forwardId int) {
	forwardIdToForwardAcceptLock.Lock()
	tcpAccept := forwardIdToForwardAccept[forwardId]
	forwardIdToForwardAcceptLock.Unlock()
	if tcpAccept != nil {
		tcpAccept.shutdown()
	}
}
