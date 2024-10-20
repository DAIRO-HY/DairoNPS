package ForwardStatisticsUtil

import (
	"DairoNPS/dao/ChannelDao"
	"DairoNPS/dao/ChannelDataStatisticsDao"
	"DairoNPS/dao/ClientDao"
	"DairoNPS/dao/ForwardDao"
	"DairoNPS/dao/SystemConfigDao"
	"DairoNPS/dao/dto"
	"sync"
	"time"
)

// 端口流量总和
var ForwardDataSizeMap = make(map[int]*ForwardDataSize)

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
		save()
		lock.Unlock()
	}
}

// 加载统计数据
func InitForwardDataSizeMap() {
	lock.Lock()
	forwardList := ForwardDao.SelectAll()

	//隧道ID对应的隧道信息
	forwardDtoMap := make(map[int]*dto.ForwardDto)
	for _, forwardDto := range forwardList {
		if forwardDto.EnableState == 0 {
			continue
		}
		forwardDtoMap[forwardDto.Id] = forwardDto
	}
	for _, forwardDto := range forwardDtoMap {
		if ForwardDataSizeMap[forwardDto.Id] != nil { //该隧道已经在统计
			continue
		}

		//加入到隧道流量统计
		ForwardDataSizeMap[forwardDto.Id] = &ForwardDataSize{
			InData:     forwardDto.InData,
			PreInData:  forwardDto.InData,
			OutData:    forwardDto.OutData,
			PreOutData: forwardDto.OutData,
		}
	}

	//移除不需要统计的对象之前,先保存统计信息
	save()

	//移除不需要统计的对象(这些对象可能已经被删除或者禁用)
	for forwardId := range ForwardDataSizeMap {
		if forwardDtoMap[forwardId] == nil {
			delete(ForwardDataSizeMap, forwardId)
		}
	}
	lock.Unlock()
}

// 保存流量记录
func save() {
	clientMap := make(map[int]*dto.ClientDto)
	for channelId, dataSize := range ForwardDataSizeMap {

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

		//添加一条统计记录
		ChannelDataStatisticsDao.Add(channelId, currentInData, currentOutData)

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
