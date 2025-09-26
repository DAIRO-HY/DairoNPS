package common

				import (
					"DairoNPS/DebugTimer"
				)

import (
	"DairoNPS/dao/ChannelDao"
	"DairoNPS/dao/ClientDao"
	"DairoNPS/web/controller/common/form"
	"net/http"
)

// post:/common/dropdown
func Dropdown(request *http.Request) map[string][]form.DropdownOutForm {
DebugTimer.Add536()

	//返回结果
	result := make(map[string][]form.DropdownOutForm)
	query := request.URL.Query()

	//组装客户端的dropdown数据------------------------------------------------------------------------------
	if query.Get("client") != "" {
DebugTimer.Add537()
		clientFormList := []form.DropdownOutForm{}
		clientList := ClientDao.SelectAll()
		for _, item := range clientList {
DebugTimer.Add538()
			if item.EnableState == 0 {
DebugTimer.Add539()
				continue
			}
			clientFormList = append(clientFormList, form.DropdownOutForm{
				Label: item.Name,
				Value: item.Id,
			})
		}
		result["client"] = clientFormList
	}

	//组装隧道的dropdown数据------------------------------------------------------------------------------
	if query.Get("channel") != "" {
DebugTimer.Add540()
		channelFormList := []form.DropdownOutForm{}
		channelList := ChannelDao.SelectAll()
		for _, item := range channelList {
DebugTimer.Add541()
			if item.EnableState == 0 {
DebugTimer.Add542()
				continue
			}
			channelFormList = append(channelFormList, form.DropdownOutForm{
				Label: item.Name,
				Value: item.Id,
			})
		}
		result["channel"] = channelFormList
	}
	return result
}
