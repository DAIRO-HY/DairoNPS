package client

import (
	"DairoNPS/dao/ChannelDataStatisticsDao"
	"DairoNPS/dao/dto"
	"DairoNPS/web"
	"DairoNPS/web/controller/data_size_log/form"
	"net/http"
	"time"
)

// 初始化
func init() {
	http.HandleFunc("/data_size/get_data_size", web.ApiHandler(getDataSize))
}

/**
 * 页面初始化数据
 * @param targetId 隧道id获取转发id
 * @param type 统计类型, 1:隧道  2:数据转发
 * @param start 开始时间
 * @param end 结束时间
 */
func getDataSize(inForm form.GetDataInForm) form.GetDataOutForm {

	//时间间隔
	timeJg := inForm.EndTime - inForm.StartTime

	//统计时间最小单位长度
	var labelFormat string
	if timeJg <= 60 { //小于1分钟，则时间最小单位到秒（yyyyMMddHHmmss）
		labelFormat = "2006-01-02 15:04:05"
	} else if timeJg <= 60*60 { //小于1小时，则时间最小单位到分（yyyyMMddHHmm）
		labelFormat = "2006-01-02 15:04"
	} else if timeJg <= 24*60*60 { //小于1天，则时间最小单位到小时（yyyyMMddHH）
		labelFormat = "2006-01-02 15"
	} else if timeJg <= 31*24*60*60 { //小于31天，则时间最小单位到天（yyyyMMdd）
		labelFormat = "2006-01-02"
	} else if timeJg <= 366*24*60*60 { //小于1一年，则时间最小单位到月（yyyyMM）
		labelFormat = "2006-01"
	} else {
		labelFormat = "2006"
	}

	label2DataForm := make(map[string]*dto.ChannelDataStatisticsDto)

	//单位
	//unit := ""
	list := ChannelDataStatisticsDao.SelectList(inForm.ClientId, inForm.ChannelId, inForm.StartTime, inForm.EndTime)
	for _, item := range list { //为每一个时间点去匹配数据
		label := time.Unix(item.Date, 0).Format(labelFormat)
		dataForm := label2DataForm[label]
		if dataForm == nil {
			dataForm = &dto.ChannelDataStatisticsDto{
				InData:  item.InData,
				OutData: item.OutData,
			}
			label2DataForm[label] = dataForm
		} else {
			dataForm.InData += item.InData
			dataForm.OutData += item.OutData
		}
	}

	//报表标题列表
	labels := []string{}

	//入网数据列表
	inDatas := []int64{}

	//出网数据列表
	outDatas := []int64{}

	loopTime := time.Unix(inForm.StartTime, 0)
	endTime := time.Unix(inForm.EndTime+1, 0)

	for loopTime.Before(endTime) {
		label := loopTime.Format(labelFormat)
		labels = append(labels, label)
		dataSize := label2DataForm[label]
		if dataSize == nil {
			inDatas = append(inDatas, 0)
			outDatas = append(outDatas, 0)
		} else {
			inDatas = append(inDatas, dataSize.InData)
			outDatas = append(outDatas, dataSize.OutData)
		}
		if labelFormat == "2006-01-02 15:04:05" { //精确到秒
			loopTime = loopTime.Add(1 * time.Second)
		} else if labelFormat == "2006-01-02 15:04" { //精确到分
			loopTime = loopTime.Add(1 * time.Minute)
		} else if labelFormat == "2006-01-02 15" { //精确到小时
			loopTime = loopTime.Add(1 * time.Hour)
		} else if labelFormat == "2006-01-02" { //精确到天
			loopTime = loopTime.AddDate(0, 0, 1)
		} else if labelFormat == "2006-01" { //精确到天
			loopTime = loopTime.AddDate(0, 1, 0)
		} else if labelFormat == "2006" { //精确到年
			loopTime = loopTime.AddDate(1, 0, 0)
		} else {
		}
	}

	return form.GetDataOutForm{
		//统计表标题列表
		Lables: labels,

		//入网流量
		InDatas: inDatas,

		//出网流量
		OutDatas: outDatas,

		// 单位
		Unit: "KB",
	}
}
