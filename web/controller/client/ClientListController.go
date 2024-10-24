package client

import (
	"DairoNPS/dao/ChannelDao"
	"DairoNPS/dao/ClientDao"
	"DairoNPS/extension/Number"
	"DairoNPS/nps/nps_client"
	"DairoNPS/web"
	"DairoNPS/web/controller/client/form"
	"net/http"
)

// 路由设置
func init() {
	http.HandleFunc("/client_list/init", web.ApiHandler(List))
	http.HandleFunc("/client_list/delete", web.ApiHandler(Delete))
	http.HandleFunc("/client_list/set_state", web.ApiHandler(SetState))
}

// 客户端列表
func List() any {
	var forms []form.ClientListForm
	clientList := ClientDao.SelectAll()
	for _, dto := range clientList {
		frm := form.ClientListForm{
			Id:          dto.Id,
			Name:        dto.Name,
			Key:         dto.Key,
			Ip:          dto.Ip,
			Version:     dto.Version,
			EnableState: dto.EnableState,
			InData:      Number.ToDataSize(dto.InData),
			OutData:     Number.ToDataSize(dto.OutData),
			IsOnline:    nps_client.IsOnline(dto.Id),
		}
		forms = append(forms, frm)
	}
	return forms
}

// 删除的表单
type DeleteForm struct {
	Id int
}

// 通过id删除一个客户端
func Delete(form DeleteForm) {

	//关闭代理监听
	nps_client.Shutdown(form.Id)
	ClientDao.Delete(form.Id)
	ChannelDao.DeleteByClient(form.Id)
}

// 修改可用状态
type SetStateForm struct {
	Id int
}

// 修改可用状态
func SetState(form SetStateForm) {
	clientDto := ClientDao.SelectOne(form.Id)
	if clientDto.EnableState == 0 {
		ClientDao.SetEnableState(form.Id, 1)
	} else {
		ClientDao.SetEnableState(form.Id, 0)
		nps_client.Shutdown(form.Id)
	}
}
