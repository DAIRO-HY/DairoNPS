package ForwardStatisticsUtil

import (
	"DairoNPS/dao/DateDataSizeDao"
	"DairoNPS/dao/ForwardDao"
	"DairoNPS/dao/SystemConfigDao"
	"DairoNPS/dao/dto"
	"sync"
	"time"
)

// 端口流量总和
var forwardDataSizeMap = make(map[int]*ForwardDataSize)

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

// 通过隧道ID获取一个统计数据
func Get(forwardId int) *ForwardDataSize {
	lock.Lock()
	dataSize := forwardDataSizeMap[forwardId]
	lock.Unlock()
	return dataSize
}

// 加载统计数据
func Init() {
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
		if forwardDataSizeMap[forwardDto.Id] != nil { //该隧道已经在统计
			continue
		}

		//加入到隧道流量统计
		forwardDataSizeMap[forwardDto.Id] = &ForwardDataSize{
			InData:     forwardDto.InData,
			PreInData:  forwardDto.InData,
			OutData:    forwardDto.OutData,
			PreOutData: forwardDto.OutData,
		}
	}

	//移除不需要统计的对象之前,先保存统计信息
	save()

	//移除不需要统计的对象(这些对象可能已经被删除或者禁用)
	for forwardId := range forwardDataSizeMap {
		if forwardDtoMap[forwardId] == nil {
			delete(forwardDataSizeMap, forwardId)
		}
	}
	lock.Unlock()
}

// 保存流量记录
func save() {
	var systemInData int64 = 0
	var systemOutData int64 = 0
	for forwardId, dataSize := range forwardDataSizeMap {

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
		DateDataSizeDao.Add(0, forwardId, currentInData, currentOutData)

		//更新端口转发出入网流量
		ForwardDao.SetDataSize(forwardId, inData, outData)

		//系统
		systemInData += currentInData
		systemOutData += currentOutData
	}

	//统计系统总流量
	SystemConfigDao.AddDataSize(systemInData, systemOutData)
}

// 获取当前统计流量总和
func GetTotal(forwardId int) (int64, int64) {
	var inData int64 = 0
	var outData int64 = 0
	lock.Lock()
	for key, dataSize := range forwardDataSizeMap {
		if forwardId != 0 { //统计某个隧道
			if key == forwardId {
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
