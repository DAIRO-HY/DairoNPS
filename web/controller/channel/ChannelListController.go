package channel

				import (
					"DairoNPS/DebugTimer"
				)

import (
	"DairoNPS/dao/ChannelDao"
	"DairoNPS/dao/ClientDao"
	"DairoNPS/dao/DateDataSizeDao"
	"DairoNPS/extension/Bool"
	"DairoNPS/extension/Number"
	"DairoNPS/nps/nps_client/tcp_client"
	"DairoNPS/nps/nps_proxy/tcp_proxy"
	"DairoNPS/nps/nps_proxy/udp_proxy"
	"DairoNPS/web/controller/channel/form"
)

// get:/channel_list
// templates:channel_list.html
func InitList() {
DebugTimer.Add508()
}

// List 隧道列表
// post:/channel_list/list
func List(clientId int) any {
DebugTimer.Add509()

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

	client := ClientDao.SelectOne(clientId)
	channelDtoList := ChannelDao.SelectByClientId(clientId)

	//返回的form表单列表
	outFormList := make([]form.ChannelListForm, len(channelDtoList))
	for i, it := range channelDtoList {
DebugTimer.Add510()
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

// Delete 通过id删除一个隧道
// post:/channel_list/delete
func Delete(id int) {
DebugTimer.Add511()

	//关闭代理监听
	tcp_proxy.ShutdownByChannel(id)
	udp_proxy.ShutdownByChannel(id)
	DateDataSizeDao.DeleteByChannelId(id)
	ChannelDao.Delete(id)
}

// SetState 修改可用状态
// post:/channel_list/set_state
func SetState(id int) {
DebugTimer.Add512()
	channel := ChannelDao.SelectOne(id)
	if channel.EnableState == 0 {
DebugTimer.Add513()
		ChannelDao.SetEnableState(id, 1)
		clientDto := ClientDao.SelectOne(channel.ClientId)
		if tcp_client.IsOnline(clientDto.Id) {
DebugTimer.Add514()
			tcp_proxy.AcceptClient(clientDto) //重新开启监听该客户端
			udp_proxy.AcceptClient(clientDto) //重新开启监听该客户端
		}
	} else {
		ChannelDao.SetEnableState(id, 0)

		//关闭代理监听
		tcp_proxy.ShutdownByChannel(channel.Id)
		udp_proxy.ShutdownByChannel(channel.Id)
	}
}
