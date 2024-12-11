package client

import (
	"DairoNPS/dao/ChannelDao"
	"DairoNPS/dao/ClientDao"
	"DairoNPS/dao/DateDataSizeDao"
	"DairoNPS/extension/Number"
	"DairoNPS/nps/nps_client/tcp_client"
	"DairoNPS/web/controller/client/form"
)

// get:/client_list
// templates:client_list.html
func InitList() {
}

// List 客户端列表
// post:/client_list/init
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

// Delete 通过id删除一个客户端
// post:/client_list/delete
func Delete(id int) {

	//关闭代理监听
	tcp_client.Shutdown(id)
	DateDataSizeDao.DeleteByClientId(id)
	ClientDao.Delete(id)
	ChannelDao.DeleteByClient(id)
}

// SetState 修改可用状态
// post:/client_list/set_state
func SetState(id int) {
	clientDto := ClientDao.SelectOne(id)
	if clientDto.EnableState == 0 {
		ClientDao.SetEnableState(id, 1)
	} else {
		ClientDao.SetEnableState(id, 0)
		tcp_client.Shutdown(id)
	}
}
