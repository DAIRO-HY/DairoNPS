package common

import (
	"DairoNPS/dao/ChannelDao"
	"DairoNPS/dao/ClientDao"
	"DairoNPS/web"
	"DairoNPS/web/controller/common/form"
	"net/http"
)

// 初始化
func init() {
	http.HandleFunc("/common/dropdown", web.ApiHandler(dropdown))
}

func dropdown(request *http.Request) map[string][]form.DropdownOutForm {

	//返回结果
	result := make(map[string][]form.DropdownOutForm)
	query := request.URL.Query()

	//组装客户端的dropdown数据------------------------------------------------------------------------------
	if query.Get("client") != "" {
		clientFormList := []form.DropdownOutForm{}
		clientList := ClientDao.SelectAll()
		for _, item := range clientList {
			if item.EnableState == 0 {
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
		channelFormList := []form.DropdownOutForm{}
		channelList := ChannelDao.SelectAll()
		for _, item := range channelList {
			if item.EnableState == 0 {
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
