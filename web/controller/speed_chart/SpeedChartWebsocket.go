package speed_chart

import (
	"DairoNPS/util/StatisticsUtil"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
)

func init() {
	http.HandleFunc("/ws/speed_chart", currentData)
}

// 创建WebSocket升级器
var upgrader = websocket.Upgrader{
	// 允许跨域请求
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocket处理函数
func currentData(w http.ResponseWriter, r *http.Request) {
	// 将HTTP连接升级为WebSocket连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("升级为WebSocket失败:", err)
		return
	}
	defer conn.Close()

	for {
		// 读取消息
		_, idData, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("读取消息失败:", err)
			break
		}
		idStr := string(idData)
		id, _ := strconv.ParseInt(idStr, 10, 64)
		channelId := int(id)
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
	}
}

func data() {

}
