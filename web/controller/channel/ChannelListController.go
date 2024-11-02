package channel

import (
	"DairoNPS/dao/ChannelDao"
	"DairoNPS/dao/ClientDao"
	"DairoNPS/dao/DateDataSizeDao"
	"DairoNPS/extension/Bool"
	"DairoNPS/extension/Number"
	"DairoNPS/nps/nps_client/tcp_client"
	"DairoNPS/nps/nps_proxy/tcp_proxy"
	"DairoNPS/nps/nps_proxy/udp_proxy"
	"DairoNPS/web"
	"DairoNPS/web/controller/channel/form"
	"net/http"
)

// 路由设置
func init() {
	http.HandleFunc("/channel_list/list", web.ApiHandler(List))
	http.HandleFunc("/channel_list/delete", web.ApiHandler(Delete))
	http.HandleFunc("/channel_list/set_state", web.ApiHandler(setState))
}

type ListInForm struct {
	ClientId int
}

// 隧道列表
func List(inForm ListInForm) any {

	////客户端下拉框数据
	//clients := ClientDao.SelectAll()
	//var clientDropdownList []form.ClientDropdownForm
	//for _, it := range clients {
	//	clientDropdownForm := form.ClientDropdownForm{
	//		Id:   it.Id,   //id
	//		Name: it.Name, //客户端名
	//	}
	//	clientDropdownList = append(clientDropdownList, clientDropdownForm)
	//}
	//
	//searchDto := dto.ChannelListSearchDto{
	//	ClientId: search.ClientId,
	//	Mode:     search.Mode,
	//}

	client := ClientDao.SelectOne(inForm.ClientId)
	channelDtoList := ChannelDao.SelectByClientId(inForm.ClientId)

	//返回的form表单列表
	outFormList := make([]form.ChannelListForm, len(channelDtoList))
	for i, it := range channelDtoList {
		outFormList[i] = form.ChannelListForm{
			Id:         it.Id,
			ClientId:   it.ClientId,
			ClientName: client.Name,
			Name:       it.Name,
			Mode:       Bool.Is(it.Mode == 1, "TCP", "UDP"),
			ServerPort: it.ServerPort,
			TargetPort: it.TargetPort,
			//EnableStateText:it.EnableStateText,
			EnableState:   it.EnableState,
			InData:        Number.ToDataSize(it.InData),
			OutData:       Number.ToDataSize(it.OutData),
			SecurityState: it.SecurityState,
			Error:         it.Error,
		}
	}
	return outFormList
}

// 删除的表单
type DeleteForm struct {
	Id int
}

// 通过id删除一个隧道
func Delete(inForm DeleteForm) {

	//关闭代理监听
	tcp_proxy.ShutdownByChannel(inForm.Id)
	udp_proxy.ShutdownByChannel(inForm.Id)
	DateDataSizeDao.DeleteByChannelId(inForm.Id)
	ChannelDao.Delete(inForm.Id)
}

// 删除的表单
type SetStateForm struct {
	Id int
}

// 修改可用状态
func setState(inForm SetStateForm) {
	channel := ChannelDao.SelectOne(inForm.Id)
	if channel.EnableState == 0 {
		ChannelDao.SetEnableState(inForm.Id, 1)
		clientDto := ClientDao.SelectOne(channel.ClientId)
		if tcp_client.IsOnline(clientDto.Id) {
			tcp_proxy.AcceptClient(clientDto) //重新开启监听该客户端
			udp_proxy.AcceptClient(clientDto) //重新开启监听该客户端
		}
	} else {
		ChannelDao.SetEnableState(inForm.Id, 0)

		//关闭代理监听
		tcp_proxy.ShutdownByChannel(channel.Id)
		udp_proxy.ShutdownByChannel(channel.Id)
	}
}
