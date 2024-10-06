package TCPPool

import (
	"net"
	"sync"
)

//客户端连接池列表

type TCPPool struct {
	ClientID int

	Socket net.Conn

	/**
	 * 标记是否已经关闭
	 */
	isClosed bool

	/**
	 * 标记是否被使用
	 */
	isUsed bool
}

var (
	lock sync.Mutex
)

/**
 * 一段时候后关闭连接池
 * 由于连接状态的不确定性,长时间无任何操作可能导致与客户端无法通讯,所以设置每个连接池在一定时间内自动销毁
 */
func CloseJob() {

	//一段时间类,如果连接池未被使用,则关闭
	//delay(CLSConfig.RECYLE_POOL_TIME)
	//lock.withLock {
	//    if (!isUsed) {//如果该连接池没有被使用则关闭它
	//        socket.close()
	//        isClosed = true
	//        TCPPoolManager.timeOutRemove(clientID, this@TCPPool)
	//    }
	//}
}

/**
 * 获取连接
 */
func GetSocket(pool TCPPool) net.Conn {
	lock.Lock()
	if !pool.isClosed { //如果连接池没有被关闭,则使用
		pool.isUsed = true

		//取消关闭等待
		//pool.closeJob.cancel()
		lock.Unlock()
		return pool.Socket
	} else {
		lock.Unlock()
		return nil
	}
}

/**
 * 关闭连接池
 */
func Close(pool TCPPool) {
	lock.Lock()
	if pool.isClosed || pool.isUsed {
		lock.Unlock()
		return
	}

	//关闭连接池
	pool.Socket.Close()

	//取消关闭等待
	//closeJob.cancel()
	pool.isClosed = true
	lock.Unlock()
}

/**
 * 发送心跳数据
 */
func sendUrgentData(pool TCPPool) {
	//pool.socket.sendUrgentData(0xFF)
}
