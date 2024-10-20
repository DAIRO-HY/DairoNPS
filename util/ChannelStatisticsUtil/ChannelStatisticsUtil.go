package ChannelStatisticsUtil

import (
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
	go timer()
}

// 计时统计
func timer() {
	for {
		time.Sleep(10 * time.Second)
		lock.Lock()
		saveStatistics()
		lock.Unlock()
	}
}

// 通过隧道ID获取一个统计数据
func Get(channelId int) *ChannelDataSize {
	lock.Lock()
	dataSize := channelDataSizeMap[channelId]
	lock.Unlock()
	return dataSize
}

// 初始化隧道统计数据
func Init() {
	lock.Lock()
	channelList := ChannelDao.SelectAll()

	//隧道ID对应的隧道信息
	channelMap := make(map[int]*dto.ChannelDto)
	for _, channel := range channelList {
		if channel.EnableState == 0 {
			continue
		}
		channelMap[channel.Id] = channel
	}
	for _, channel := range channelMap {
		if channelDataSizeMap[channel.Id] != nil { //该隧道已经在统计
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
		if channelMap[channelId] == nil {
			delete(channelDataSizeMap, channelId)
		}
	}
	lock.Unlock()
}

// 保存流量记录
func saveStatistics() {
	clientMap := make(map[int]*dto.ClientDto)
	for channelId, dataSize := range channelDataSizeMap {

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
		inData += client.InData
		outData += client.OutData
		ClientDao.SetDataSize(clientId, client.InData, client.OutData)
	}

	//统计系统总流量
	SystemConfigDao.AddDataSize(inData, outData)
}

// 获取当前统计流量总和
func GetTotal(clientId int, channelId int) (int64, int64) {
	var inData int64 = 0
	var outData int64 = 0
	lock.Lock()
	for cid, dataSize := range channelDataSizeMap {
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
	lock.Unlock()
	return inData, outData
}
