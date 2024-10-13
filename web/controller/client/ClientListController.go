package client

import (
	"DairoNPS/dao/ClientDao"
	"DairoNPS/web"
	"DairoNPS/web/controller/client/form"
	"net/http"
)

// 路由设置
func init() {
	http.HandleFunc("/client_list/init", web.ApiHandler(List))
}

// 客户端列表
func List() any {
	var forms []form.ClientListForm
	clientList := ClientDao.SelectAll()
	for _, dto := range clientList {
		frm := form.ClientListForm{
			Id:           dto.Id,
			Name:         dto.Name,
			Key:          dto.Key,
			Ip:           dto.Ip,
			Version:      dto.Version,
			EnableState:  dto.EnableState,
			InDataTotal:  dto.InDataTotal,
			OutDataTotal: dto.OutDataTotal,
			OnlineState:  true,
			//onlineState: CLSStatus.isClientOnline(it.id!!)
		}
		forms = append(forms, frm)
	}
	return forms
}
