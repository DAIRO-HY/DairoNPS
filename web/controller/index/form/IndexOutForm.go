package form

type IndexOutForm struct {

	// 正在监听的隧道
	ProxyCount int

	// 在线客户端数量
	OnlineClientCount int

	// 桥接通信中的数量
	TcpBridgeCount int

	// 当前TCP连接池
	TcpPoolCount int

	// 当前UDP会话数
	UdpSessionCount int

	// 当前UDP连接池
	UdpPoolCount int

	// 入网流量
	InDataTotal string

	// 出网流量
	OutDataTotal string

	// 当前正在代理服务数
	ForwardCount int

	// 代理服务会话数
	ForwardBridgeCount int

	NumGoroutine int    //当前协程数
	Memory       string //内存分配
	//SystemMemory string //系统内存占用
	//HeapAlloc    string //堆内存分配
	//HeapSys      string //堆内存系统占用
	//NumGC        uint32 //垃圾回收次数
}
