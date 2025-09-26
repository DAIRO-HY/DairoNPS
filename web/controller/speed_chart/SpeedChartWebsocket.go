package speed_chart

				import (
					"DairoNPS/DebugTimer"
				)

import (
	"DairoNPS/util/ChannelStatisticsUtil"
	"DairoNPS/util/ForwardStatisticsUtil"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
	"strings"
)

// CurrentData WebSocket处理函数
// post:/ws/speed_chart
func CurrentData(w http.ResponseWriter, r *http.Request) {
DebugTimer.Add594()
	// 创建WebSocket升级器
	var upgrader = websocket.Upgrader{
		// 允许跨域请求
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	// 将HTTP连接升级为WebSocket连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
DebugTimer.Add595()
		fmt.Println("升级为WebSocket失败:", err)
		return
	}
	defer conn.Close()

	for {
DebugTimer.Add596()
		// 读取消息
		_, idData, err := conn.ReadMessage()
		if err != nil {
DebugTimer.Add597()
			break
		}
		clientId := 0
		channelId := 0
		forwardId := 0
		idStr := string(idData)
		id, _ := strconv.ParseInt(idStr[1:], 10, 64)
		if strings.HasPrefix(idStr, "C") { //获取客户端ID
DebugTimer.Add598()
			clientId = int(id)
		} else if strings.HasPrefix(idStr, "N") { //获取隧道ID
			channelId = int(id)
		} else if strings.HasPrefix(idStr, "F") { //获取端口转发ID
			forwardId = int(id)
		} else {
		}
		var channelInData int64
		var channelOutData int64
		var forwardInData int64
		var forwardOutData int64
		if clientId == 0 && channelId == 0 && forwardId == 0 { //统计所有
DebugTimer.Add599()

			//隧道流量总和
			channelInData, channelOutData = ChannelStatisticsUtil.GetTotal(0, 0)

			//端口转发流量总和
			forwardInData, forwardOutData = ForwardStatisticsUtil.GetTotal(0)
		} else if clientId != 0 { //客户端流量总和
			channelInData, channelOutData = ChannelStatisticsUtil.GetTotal(clientId, 0)
		} else if channelId != 0 { //隧道流量总和
			channelInData, channelOutData = ChannelStatisticsUtil.GetTotal(0, channelId)
		} else if forwardId != 0 { //客户端流量总和
			forwardInData, forwardOutData = ForwardStatisticsUtil.GetTotal(forwardId)
		} else {
		}

		message := strconv.FormatInt(channelInData+forwardInData, 10) + ":" + strconv.FormatInt(channelOutData+forwardOutData, 10)

		// 发送消息
		err = conn.WriteMessage(websocket.TextMessage, []uint8(message))
		if err != nil {
DebugTimer.Add600()
			break
		}
	}
}
