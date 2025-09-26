package channel

				import (
					"DairoNPS/DebugTimer"
				)

import (
	"DairoNPS/dao/ChannelDao"
	"DairoNPS/dao/ClientDao"
	"DairoNPS/dao/ForwardDao"
	"DairoNPS/dao/dto"
	"DairoNPS/extension/Bool"
	"DairoNPS/extension/Date"
	"DairoNPS/extension/Number"
	"DairoNPS/nps/nps_client/tcp_client"
	"DairoNPS/nps/nps_proxy/tcp_proxy"
	"DairoNPS/nps/nps_proxy/udp_proxy"
	"DairoNPS/web/controller"
	"DairoNPS/web/controller/channel/form"
	"fmt"
)

// get:/channel_list/channel_edit
// templates:channel_edit.html
func InitEdit() {
DebugTimer.Add493()
}

// Info 隧道编辑
// post:/channel_list/channel_edit/info
func Info(ClientId int, Id int) any {
DebugTimer.Add494()
	client := ClientDao.SelectOne(ClientId)
	var outForm form.ChannelEditForm
	if Id == 0 {
DebugTimer.Add495()
		outForm = form.ChannelEditForm{
			Mode: 1,
		}
	} else { //修改时
		channelDto := ChannelDao.SelectOne(Id)
		outForm = form.ChannelEditForm{
			Id:            channelDto.Id,
			Name:          channelDto.Name,
			Mode:          channelDto.Mode,
			Remark:        channelDto.Remark,
			ServerPort:    channelDto.ServerPort,
			TargetPort:    channelDto.TargetPort,
			Date:          Date.FormatByTimespan(channelDto.Date),
			EnableState:   Bool.Is(channelDto.EnableState == 0, "关闭", "开启"),
			InData:        Number.ToDataSize(channelDto.InData),
			OutData:       Number.ToDataSize(channelDto.OutData),
			SecurityState: channelDto.SecurityState,
		}
	}
	outForm.ClientId = ClientId
	outForm.ClientName = client.Name
	return outForm
}

// Edit 提交表单API
// post:/channel_list/channel_edit/edit
func Edit(form form.ChannelEditForm) any {
DebugTimer.Add496()
	//try {
	err := validate(form)
	if err != nil {
DebugTimer.Add497()
		return err
	}
	//val dto = if (form.id != null) {//更新时
	//    ChannelDao.selectOne(form.id!!) ?: throw BusinessException("该客户端不存在")
	//} else {
	//    ChannelDto()
	//}

	channel := &dto.ChannelDto{
		Id:            form.Id,
		Name:          form.Name,
		ServerPort:    form.ServerPort,
		Mode:          form.Mode,
		ClientId:      form.ClientId,
		TargetPort:    form.TargetPort,
		SecurityState: form.SecurityState,
		AclState:      form.AclState,
		Remark:        form.Remark,
	}
	if form.Id == 0 {
DebugTimer.Add498()
		ChannelDao.Add(channel)
	} else { //更新时
		ChannelDao.Update(channel)
	}

	//关闭代理监听
	tcp_proxy.ShutdownByChannel(channel.Id)
	udp_proxy.ShutdownByChannel(channel.Id)
	clientDto := ClientDao.SelectOne(channel.ClientId)
	if tcp_client.IsOnline(clientDto.Id) {
DebugTimer.Add499()
		tcp_proxy.AcceptClient(clientDto) //重新开启监听该客户端
		udp_proxy.AcceptClient(clientDto) //重新开启监听该客户端
	}
	return nil
}

// 表单验证
func validate(inForm form.ChannelEditForm) error {
DebugTimer.Add500()
	if len(inForm.Name) == 0 {
DebugTimer.Add501()
		return &controller.BusinessException{
			Message: "请填写隧道名",
		}
	}
	if len(inForm.Name) > 32 {
DebugTimer.Add502()
		return &controller.BusinessException{
			Message: "隧道名长度不能超过32个字符",
		}
	}
	//if len(inForm.ServerPort) == 0 {
	//	return &controller.BusinessException{
	//		Message: "服务端口必须设置",
	//	}
	//}
	//port, err := strconv.ParseInt(inForm.ServerPort, 10, 64)
	//if err != nil {
	//	return &controller.BusinessException{
	//		Message: "服务端口必须是一个数字",
	//	}
	//}
	if inForm.ServerPort < 0 || inForm.ServerPort > 65535 {
DebugTimer.Add503()
		return &controller.BusinessException{
			Message: "服务端口必须在0到65535之间",
		}
	}
	portChannel := ChannelDao.SelectByPort(inForm.ServerPort)
	if inForm.Id == 0 { //创建时
DebugTimer.Add504()
		if portChannel != nil {
DebugTimer.Add505()
			return &controller.BusinessException{
				Message: fmt.Sprintf("端口:%d已经被其他隧道占用", inForm.ServerPort),
			}
		}
	} else {
		if portChannel != nil && portChannel.Id != inForm.Id {
DebugTimer.Add506()
			return &controller.BusinessException{
				Message: fmt.Sprintf("端口:%d已经被其他隧道占用", inForm.ServerPort),
			}
		}
	}
	portForward := ForwardDao.SelectByPort(inForm.ServerPort)
	if portForward != nil {
DebugTimer.Add507()
		return &controller.BusinessException{
			Message: fmt.Sprintf("端口:%d已被端口转发:%s 占用", portForward.Port, portForward.Name),
		}
	}
	return nil
}
