package index

import (
	"DairoNPS/dao/SystemConfigDao"
	"DairoNPS/extension/Number"
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
)

// 初始化
func init() {
	http.HandleFunc("/index/data", web.ApiHandler(data))
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
	outForm := form.IndexOutForm{
		OnlineClientCount: nps_client.OnlineCount(),             //在线客户端数量
		TcpBridgeCount:    nps_bridge.GetBridgeCount(),          //当前TCP会话数
		TcpPoolCount:      nps_pool.GetPoolCount(),              //当前TCP连接池
		UdpSessionCount:   0,                                    //当前UDP会话数
		UdpPoolCount:      0,                                    //当前UDP连接池
		InDataTotal:       Number.ToDataSize(systemDto.InData),  //入网流量
		OutDataTotal:      Number.ToDataSize(systemDto.OutData), //出网流量
		ProxyCount:        nps_channel_proxy.GetProxyCount(),    //当前正在代理数
		//ForwardCount:       CLSStatus.forwardCount,               //当前正在代理服务数
		//ForwardBridgeCount: CLSStatus.forwardBridgeCount,         //代理服务会话数
	}
	return outForm
}
