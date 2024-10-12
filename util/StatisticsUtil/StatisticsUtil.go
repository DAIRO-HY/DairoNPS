package StatisticsUtil

import (
	"DairoNPS/dao/ChannelDao"
	"fmt"
	"time"
)

// 隧道流量总和
var ChannelDataSizeMap = make(map[int]*ChannelDataSizeLog)

func Init() {
	channelAll := ChannelDao.SelectAll()
	for i := 0; i < len(channelAll); i++ {
		channel := channelAll[i]
		if channel.EnableState == 0 {
			continue
		}
		ChannelDataSizeMap[channel.Id] = &ChannelDataSizeLog{
			InData:  channel.InDataTotal,
			OutData: channel.OutDataTotal,
		}
	}
}

// 流量统计
func ThransferDataSizeStatistics() {
	for {
		time.Sleep(5 * time.Second)
		for k, v := range ChannelDataSizeMap {
			fmt.Printf("流量统计%d=入:%d 出:%d\n", k, v.InData, v.OutData)
		}
	}
}
