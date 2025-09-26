package index

				import (
					"DairoNPS/DebugTimer"
				)

import (
	"DairoNPS/dao/SystemConfigDao"
	"DairoNPS/extension/Number"
	"DairoNPS/forward"
	"DairoNPS/nps/nps_bridge/tcp_bridge"
	"DairoNPS/nps/nps_bridge/udp_bridge"
	"DairoNPS/nps/nps_client/tcp_client"
	"DairoNPS/nps/nps_pool/tcp_pool"
	"DairoNPS/nps/nps_pool/udp_pool"
	"DairoNPS/nps/nps_proxy/tcp_proxy"
	"DairoNPS/web/controller/index/form"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"runtime"
)

// Home 首页重定向
// get:/
// templates:index.html
func Home() {
DebugTimer.Add574()
}

// get:/index
// templates:index.html
func Init() {
DebugTimer.Add575()
}

// Data WebSocket处理函数
// post:/index/data
func Data(writer http.ResponseWriter, request *http.Request) { // 创建WebSocket升级器
DebugTimer.Add576()
	var upgrader = websocket.Upgrader{
		// 允许跨域请求
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	// 将HTTP连接升级为WebSocket连接
	conn, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
DebugTimer.Add577()
		fmt.Println("升级为WebSocket失败:", err)
		return
	}
	defer conn.Close()

	for {
DebugTimer.Add578()
		// 读取消息
		_, _, err := conn.ReadMessage()
		if err != nil {
DebugTimer.Add579()
			break
		}
		data := getData()
		jsonData, err := json.Marshal(data)
		if err != nil {
DebugTimer.Add580()
			break
		}

		// 发送消息
		err = conn.WriteMessage(websocket.TextMessage, jsonData)
		if err != nil {
DebugTimer.Add581()
			break
		}
	}
}

// Gc 垃圾回收
// post:/index/gc
func Gc() {
DebugTimer.Add582()
	runtime.GC()
}

// 页面初始化
func getData() form.IndexOutForm {
DebugTimer.Add583()
	systemDto := SystemConfigDao.SelectOne()

	// 获取内存使用情况
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	outForm := form.IndexOutForm{
		OnlineClientCount:  tcp_client.OnlineCount(),             //在线客户端数量
		TcpBridgeCount:     tcp_bridge.GetBridgeCount(),          //当前TCP桥接数
		TcpPoolCount:       tcp_pool.GetPoolCount(),              //当前TCP连接池
		UdpBridgeCount:     udp_bridge.GetBridgeCount(),          //当前UDP桥接数
		UdpPoolCount:       udp_pool.GetPoolCount(),              //当前UDP连接池
		InDataTotal:        Number.ToDataSize(systemDto.InData),  //入网流量
		OutDataTotal:       Number.ToDataSize(systemDto.OutData), //出网流量
		ProxyCount:         tcp_proxy.GetProxyCount(),            //当前正在代理数
		ForwardCount:       forward.GetAcceptCount(),             //端口转发代理端口数量
		ForwardBridgeCount: forward.GetBridgeCount(),             //端口转发当前桥接数量
		NumGoroutine:       runtime.NumGoroutine(),               //当前协程数
		Memory:             Number.ToDataSize(memStats.Alloc),    //内存分配
		//SystemMemory:       Number.ToDataSize(memStats.Sys),       //系统内存占用
		//HeapAlloc:          Number.ToDataSize(memStats.HeapAlloc), //堆内存分配
		//HeapSys:            Number.ToDataSize(memStats.HeapSys),   //堆内存系统占用
		NumGC: memStats.NumGC, //垃圾回收次数
	}
	return outForm
}
