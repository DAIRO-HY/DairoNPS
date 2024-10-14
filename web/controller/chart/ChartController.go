package chart

import (
	"DairoNPS/util/StatisticsUtil"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

func init() {
	http.HandleFunc("/ws", wsHandler)
}

// 创建WebSocket升级器
var upgrader = websocket.Upgrader{
	// 允许跨域请求
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocket处理函数
func wsHandler(w http.ResponseWriter, r *http.Request) {
	// 将HTTP连接升级为WebSocket连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("升级为WebSocket失败:", err)
		return
	}
	defer conn.Close()

	for {
		// 读取消息
		//mt, message, err := conn.ReadMessage()
		//if err != nil {
		//	fmt.Println("读取消息失败:", err)
		//	break
		//}
		//fmt.Printf("收到消息: %s\n", message)

		channelId := 1
		channelData := StatisticsUtil.ChannelDataSizeMap[channelId]

		var message string
		if channelData == nil {
			message = "0:0"
		} else {
			message = fmt.Sprintf("%d:%d", channelData.InData, channelData.OutData)
		}

		// 发送消息
		err = conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			fmt.Println("发送消息失败:", err)
			break
		}
		time.Sleep(1 * time.Second)
	}
}

func data() {

}
