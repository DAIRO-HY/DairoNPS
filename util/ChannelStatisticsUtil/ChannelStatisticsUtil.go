package ChannelStatisticsUtil

				import (
					"DairoNPS/DebugTimer"
				)

import (
	"DairoNPS/constant/NPSConstant"
	"DairoNPS/dao/ChannelDao"
	"DairoNPS/dao/ClientDao"
	"DairoNPS/dao/DateDataSizeDao"
	"DairoNPS/dao/SystemConfigDao"
	"DairoNPS/dao/dto"
	"sync"
	"time"
)

// 隧道流量总和
var channelDataSizeMap = make(map[int]*ChannelDataSize)

// 统计锁
var lock sync.Mutex

// 流量统计
func init() {
DebugTimer.Add376()
	go timer()
}

// 计时统计
func timer() {
DebugTimer.Add377()
	for {
DebugTimer.Add378()
		time.Sleep(NPSConstant.STATISTICS_DATA_SIZE_TIMER * time.Second)
		lock.Lock()
		saveStatistics()
		lock.Unlock()
	}
}

// 通过隧道ID获取一个统计数据
func Get(channelId int) *ChannelDataSize {
DebugTimer.Add379()
	lock.Lock()
	dataSize := channelDataSizeMap[channelId]
	lock.Unlock()
	return dataSize
}

// 初始化隧道统计数据
func Init() {
DebugTimer.Add380()
	lock.Lock()
	channelList := ChannelDao.SelectAll()

	//隧道ID对应的隧道信息
	channelMap := make(map[int]*dto.ChannelDto)
	for _, channel := range channelList {
DebugTimer.Add381()
		if channel.EnableState == 0 {
DebugTimer.Add382()
			continue
		}
		channelMap[channel.Id] = channel
	}
	for _, channel := range channelMap {
DebugTimer.Add383()
		if channelDataSizeMap[channel.Id] != nil { //该隧道已经在统计
DebugTimer.Add384()
			continue
		}

		//加入到隧道流量统计
		channelDataSizeMap[channel.Id] = &ChannelDataSize{
			ClientId:   channel.ClientId,
			InData:     channel.InData,
			PreInData:  channel.InData,
			OutData:    channel.OutData,
			PreOutData: channel.OutData,
		}
	}

	//移除不需要统计的隧道之前,先保存统计信息
	saveStatistics()

	//移除不需要统计的隧道(这些对象可能已经被删除或者禁用)
	for channelId := range channelDataSizeMap {
DebugTimer.Add385()
		if channelMap[channelId] == nil {
DebugTimer.Add386()
			delete(channelDataSizeMap, channelId)
		}
	}
	lock.Unlock()
}

// 保存流量记录
func saveStatistics() {
DebugTimer.Add387()
	clientMap := make(map[int]*dto.ClientDto)
	for channelId, dataSize := range channelDataSizeMap {
DebugTimer.Add388()

		//当前流量(入网)
		inData := dataSize.InData

		//上次统计到的流量(入网)
		preInData := dataSize.PreInData

		//本次统计变更(入网)
		currentInData := inData - preInData

		//当前流量(出网)
		outData := dataSize.OutData

		//上次统计到的流量(出网)
		preOutdata := dataSize.PreOutData

		//本次统计变更(出网)
		currentOutData := outData - preOutdata

		if currentInData == 0 && currentOutData == 0 { //没有数据变化时跳过
DebugTimer.Add389()
			continue
		}

		//更新本次统计(入网)
		dataSize.PreInData = inData

		//更新本次统计(出网)
		dataSize.PreOutData = outData

		//添加一条统计记录
		DateDataSizeDao.Add(channelId, 0, currentInData, currentOutData)

		//更新隧道出入网流量
		ChannelDao.SetDataSize(channelId, inData, outData)

		//统计客户端流量
		clientDto := clientMap[dataSize.ClientId]
		if clientDto == nil {
DebugTimer.Add390()
			clientMap[dataSize.ClientId] = &dto.ClientDto{
				InData:  currentInData,
				OutData: currentOutData,
			}
		} else {
			clientDto.InData += currentInData
			clientDto.OutData += currentOutData
		}
	}

	var inData int64 = 0
	var outData int64 = 0

	//统计客户端入出网流量
	for clientId, client := range clientMap {
DebugTimer.Add391()
		inData += client.InData
		outData += client.OutData
		ClientDao.SetDataSize(clientId, client.InData, client.OutData)
	}

	//统计系统总流量
	SystemConfigDao.AddDataSize(inData, outData)
}

// 获取当前统计流量总和
func GetTotal(clientId int, channelId int) (int64, int64) {
DebugTimer.Add392()
	var inData int64 = 0
	var outData int64 = 0
	lock.Lock()
	for cid, dataSize := range channelDataSizeMap {
DebugTimer.Add393()
		if channelId != 0 { //统计某个隧道
DebugTimer.Add394()
			if cid == channelId {
DebugTimer.Add395()
				inData += dataSize.InData
				outData += dataSize.OutData
			}
		} else if clientId != 0 { //统计某个客户端
			if dataSize.ClientId == clientId {
DebugTimer.Add396()
				inData += dataSize.InData
				outData += dataSize.OutData
			}
		} else { //统计所有
			inData += dataSize.InData
			outData += dataSize.OutData
		}
	}
	lock.Unlock()
	return inData, outData
}
