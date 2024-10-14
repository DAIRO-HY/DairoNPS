package channel

import (
	"DairoNPS/dao/ChannelDao"
	"DairoNPS/web"
	"DairoNPS/web/controller/channel/form"
	"net/http"
)

// 路由设置
func init() {
	http.HandleFunc("/channel_list/list", web.ApiHandler(List))
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

	channelDtoList := ChannelDao.SelectByClientId(inForm.ClientId)

	//返回的form表单列表
	outFormList := make([]form.ChannelListForm, len(channelDtoList))
	for i, it := range channelDtoList {
		outFormList[i] = form.ChannelListForm{
			Id:       it.Id,
			ClientId: it.ClientId,
			//ClientName: it.ClientName,
			Name:       it.Name,
			Mode:       it.Mode,
			ServerPort: it.ServerPort,
			TargetPort: it.TargetPort,
			//EnableStateText:it.EnableStateText,
			EnableState:   it.EnableState,
			InDataTotal:   it.InDataTotal,
			OutDataTotal:  it.OutDataTotal,
			SecurityState: it.SecurityState,
		}
	}
	return outFormList
}

///**
// * 通过id删除一个隧道
// */
//@PostMapping("/delete")
//@ResponseBody
//fun delete(id: Int) {
//GlobalScope.launch(Dispatchers.IO) {
//
////关闭隧道之后再打开
//CLServer.closeByChannel(id)
//ChannelDao.delete(id)
//}
//}
//
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
