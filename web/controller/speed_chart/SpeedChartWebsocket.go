package speed_chart

import (
	"DairoNPS/util/StatisticsUtil"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
	"strings"
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
			break
		}
		channelId := 0
		clientId := 0
		idStr := string(idData)
		id, _ := strconv.ParseInt(idStr[1:], 10, 64)
		if strings.HasPrefix(idStr, "C") {
			clientId = int(id)
		} else if strings.HasPrefix(idStr, "N") {
			channelId = int(id)
		} else {
		}
		var inData int64 = 0
		var outData int64 = 0
		StatisticsUtil.Lock.Lock()
		for cid, dataSize := range StatisticsUtil.ChannelDataSizeMap {
			if channelId != 0 { //统计某个隧道
				if cid == channelId {
					inData += dataSize.InData
					outData += dataSize.OutData
				}
			} else if clientId != 0 { //统计某个客户端
				if dataSize.ClientId == clientId {
					inData += dataSize.InData
					outData += dataSize.OutData
				}
			} else { //统计所有
				inData += dataSize.InData
				outData += dataSize.OutData
			}
		}
		StatisticsUtil.Lock.Unlock()
		message := strconv.FormatInt(inData, 10) + ":" + strconv.FormatInt(outData, 10)

		// 发送消息
		err = conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			break
		}
	}
}
