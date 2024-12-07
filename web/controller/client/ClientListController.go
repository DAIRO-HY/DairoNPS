package client

import (
	"DairoNPS/dao/ChannelDao"
	"DairoNPS/dao/ClientDao"
	"DairoNPS/dao/DateDataSizeDao"
	"DairoNPS/extension/Number"
	"DairoNPS/nps/nps_client/tcp_client"
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
			IsOnline:    tcp_client.IsOnline(dto.Id),
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
func Delete(inForm DeleteForm) {

	//关闭代理监听
	tcp_client.Shutdown(inForm.Id)
	DateDataSizeDao.DeleteByClientId(inForm.Id)
	ClientDao.Delete(inForm.Id)
	ChannelDao.DeleteByClient(inForm.Id)
}

// 修改可用状态
type SetStateForm struct {
	Id int
}

// 修改可用状态
func SetState(inForm SetStateForm) {
	clientDto := ClientDao.SelectOne(inForm.Id)
	if clientDto.EnableState == 0 {
		ClientDao.SetEnableState(inForm.Id, 1)
	} else {
		ClientDao.SetEnableState(inForm.Id, 0)
		tcp_client.Shutdown(inForm.Id)
	}
}
