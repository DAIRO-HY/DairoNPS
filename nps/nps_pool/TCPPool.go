package nps_pool

import (
	"net"
	"sync"
)

//客户端连接池列表

type TCPPool struct {
	ClientID int

	Socket net.Conn

	// 创建时间
	CreateTime int64

	// 标记是否已经关闭
	IsClosed bool

	// 标记是否被使用
	isUsed bool
}

var lock sync.Mutex

/**
 * 一段时候后关闭连接池
 * 由于连接状态的不确定性,长时间无任何操作可能导致与客户端无法通讯,所以设置每个连接池在一定时间内自动销毁
 */
func (mine *TCPPool) CloseJob() {

	//一段时间类,如果连接池未被使用,则关闭
	//delay(CLSConfig.RECYLE_POOL_TIME)
	//lock.withLock {
	//    if (!isUsed) {//如果该连接池没有被使用则关闭它
	//        socket.close()
	//        IsClosed = true
	//        TCPPoolManager.timeOutRemove(clientID, this@TCPPool)
	//    }
	//}
}

/**
 * 获取连接
 */
func (mine *TCPPool) GetSocket() net.Conn {
	lock.Lock()
	if !mine.IsClosed { //如果连接池没有被关闭,则使用
		mine.isUsed = true

		//取消关闭等待
		//pool.closeJob.cancel()
		lock.Unlock()
		return mine.Socket
	} else {
		lock.Unlock()
		return nil
	}
}

/**
 * 关闭连接池
 */
func (mine *TCPPool) Shutdown() {
	lock.Lock()
	if mine.IsClosed || mine.isUsed {
		lock.Unlock()
		return
	}

	//关闭连接池
	mine.Socket.Close()

	//取消关闭等待
	//closeJob.cancel()
	mine.IsClosed = true
	lock.Unlock()
}

/**
 * 发送心跳数据
 */
func (mine *TCPPool) sendUrgentData() {
	//mine.socket.sendUrgentData(0xFF)
}
