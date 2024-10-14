package channel

import (
	"DairoNPS/bridge"
	"DairoNPS/dao/ChannelDao"
	"DairoNPS/dao/ClientDao"
	"DairoNPS/extension/Number"
	"DairoNPS/proxy"
	"DairoNPS/web"
	"DairoNPS/web/controller/channel/form"
	"net/http"
)

// 路由设置
func init() {
	http.HandleFunc("/channel_list/list", web.ApiHandler(List))
	http.HandleFunc("/channel_list/delete", web.ApiHandler(Delete))
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
			Mode:       it.Mode,
			ServerPort: it.ServerPort,
			TargetPort: it.TargetPort,
			//EnableStateText:it.EnableStateText,
			EnableState:   it.EnableState,
			InDataTotal:   Number.ToDataSize(it.InDataTotal),
			OutDataTotal:  Number.ToDataSize(it.OutDataTotal),
			SecurityState: it.SecurityState,
		}
	}
	return outFormList
}

// 删除的表单
type DeleteForm struct {
	Id int
}

// 通过id删除一个隧道
func Delete(form DeleteForm) {

	//关闭代理监听
	proxy.CloseByChannel(form.Id)

	//关闭隧道所有正在通信的连接
	bridge.CloseByChannel(form.Id)
	ChannelDao.Delete(form.Id)
}

///**
// * 修改可用状态
// */
//@PostMapping("/set_state")
//@ResponseBody
//fun setState(id: Int) {
//val channel = ChannelDao.selectOne(id)!!
//channel.enableState = if (channel.enableState == 0) 1 else 0
//ChannelDao.update(channel)
//GlobalScope.launch(Dispatchers.IO) {
//if (channel.enableState == 0) {
//CLServer.closeByChannel(id)
//} else {
//CLServer.start(id)
//}
//}
//}
//}
