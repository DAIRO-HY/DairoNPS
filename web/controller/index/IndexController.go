package index

import (
	"DairoNPS/dao/SystemConfigDao"
	"DairoNPS/extension/Number"
	"DairoNPS/forward"
	"DairoNPS/nps/nps_bridge"
	"DairoNPS/nps/nps_channel_proxy"
	"DairoNPS/nps/nps_client"
	"DairoNPS/nps/nps_pool"
	"DairoNPS/web"
	"DairoNPS/web/controller/index/form"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"runtime"
)

// 初始化
func init() {
	http.HandleFunc("/index/data", web.ApiHandler(data))
	http.HandleFunc("/index/gc", web.ApiHandler(gc))
}

// 创建WebSocket升级器
var upgrader = websocket.Upgrader{
	// 允许跨域请求
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocket处理函数
func data(w http.ResponseWriter, r *http.Request) {
	// 将HTTP连接升级为WebSocket连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("升级为WebSocket失败:", err)
		return
	}
	defer conn.Close()

	for {
		// 读取消息
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
		data := getData()
		jsonData, err := json.Marshal(data)
		if err != nil {
			break
		}

		// 发送消息
		err = conn.WriteMessage(websocket.TextMessage, jsonData)
		if err != nil {
			break
		}
	}
}

// 页面初始化
func getData() form.IndexOutForm {
	systemDto := SystemConfigDao.SelectOne()

	// 获取内存使用情况
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	outForm := form.IndexOutForm{
		OnlineClientCount:  nps_client.OnlineCount(),             //在线客户端数量
		TcpBridgeCount:     nps_bridge.GetBridgeCount(),          //当前TCP会话数
		TcpPoolCount:       nps_pool.GetPoolCount(),              //当前TCP连接池
		UdpSessionCount:    0,                                    //当前UDP会话数
		UdpPoolCount:       0,                                    //当前UDP连接池
		InDataTotal:        Number.ToDataSize(systemDto.InData),  //入网流量
		OutDataTotal:       Number.ToDataSize(systemDto.OutData), //出网流量
		ProxyCount:         nps_channel_proxy.GetProxyCount(),    //当前正在代理数
		ForwardCount:       forward.GetAcceptCount(),             //端口转发代理端口数量
		ForwardBridgeCount: forward.GetBridgeCount(),             //端口转发当前桥接数量
		NumGoroutine:       runtime.NumGoroutine(),               //当前协程数
		Memory:             Number.ToDataSize(memStats.Alloc),    //内存分配
		//SystemMemory:       Number.ToDataSize(memStats.Sys),       //系统内存占用
		//HeapAlloc:          Number.ToDataSize(memStats.HeapAlloc), //堆内存分配
		//HeapSys:            Number.ToDataSize(memStats.HeapSys),   //堆内存系统占用
		//NumGC:              memStats.NumGC,                        //垃圾回收次数
	}
	return outForm
}

// 垃圾回收
func gc() {
	runtime.GC()
}
