package StatisticsUtil

import (
	"DairoNPS/dao/ChannelDao"
	"DairoNPS/dao/ChannelDataStatisticsDao"
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
			InData:     channel.InDataTotal,
			PreInData:  channel.InDataTotal,
			OutData:    channel.OutDataTotal,
			PreOutData: channel.OutDataTotal,
		}
	}
}

// 流量统计
func Statistics() {
	for {
		time.Sleep(1 * time.Second)
		for channelId, dataSize := range ChannelDataSizeMap {

			//当前流量(入网)
			inData := dataSize.InData

			//上次统计到的流量(入网)
			preIndata := dataSize.PreInData

			//更新本次统计(入网)
			dataSize.PreInData = inData

			//本次统计变更(入网)
			currentInData := inData - preIndata

			//当前流量(出网)
			outData := dataSize.OutData

			//上次统计到的流量(出网)
			preOutdata := dataSize.PreOutData

			//更新本次统计(出网)
			dataSize.PreOutData = outData

			//本次统计变更(出网)
			currentOutData := outData - preOutdata

			if currentInData == 0 && currentOutData == 0 { //没有数据变化时跳过
				continue
			}

			//更新隧道出入网流量
			ChannelDao.SetDataSize(channelId, inData, outData)

			//添加一条统计记录
			ChannelDataStatisticsDao.Add(channelId, currentInData, currentOutData)
		}
	}
}
