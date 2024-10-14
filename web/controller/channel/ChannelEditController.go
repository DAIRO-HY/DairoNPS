package channel

import (
	"DairoNPS/dao/ChannelDao"
	"DairoNPS/dao/ClientDao"
	"DairoNPS/dao/dto"
	"DairoNPS/extension/Bool"
	"DairoNPS/extension/Date"
	"DairoNPS/extension/Number"
	"DairoNPS/proxy"
	"DairoNPS/web"
	"DairoNPS/web/controller"
	"DairoNPS/web/controller/channel/form"
	"fmt"
	"net/http"
)

// 路由设置
func init() {
	http.HandleFunc("/channel_list/channel_edit/info", web.ApiHandler(Info))
	http.HandleFunc("/channel_list/channel_edit/edit", web.ApiHandler(Edit))
}

// 输入参数
type EditInForm struct {
	ClientId int
	Id       int
}

// 隧道编辑
func Info(inForm EditInForm) any {
	client := ClientDao.SelectOne(inForm.ClientId)
	var outForm form.ChannelEditForm
	if inForm.Id == 0 {
		outForm = form.ChannelEditForm{
			Mode: 1,
		}
	} else { //修改时
		channelDto := ChannelDao.SelectOne(inForm.Id)
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

		//val aclDtoList = ChannelAclDao.selectByChannelId(id)
		//val aclIp = aclDtoList.map{it.remark + ":" + it.ip}.toTypedArray()
		//request.setAttribute("aclIp", ObjectMapper().writeValueAsString(aclIp))
	}
	outForm.ClientId = inForm.ClientId
	outForm.ClientName = client.Name
	return outForm
}

// 提交表单API
func Edit(form form.ChannelEditForm) any {
	//try {
	err := validate(form)
	if err != nil {
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
		ChannelDao.Add(channel)
	} else { //更新时
		ChannelDao.Update(channel)
	}

	////判断是否IP地址的正则表达式
	//val ipRegex = """^(([01]?\d{1,2}|2[0-4]\d|25[0-5])\.){3}([01]?\d{1,2}|2[0-4]\d|25[0-5])$""".toRegex()
	//val aclDtoList = form.aclIp.map {
	//    val aclDto = ChannelAclDto()
	//    if (it.contains(":")) {
	//        it.split(":").also {
	//            aclDto.remark = it[0]
	//            aclDto.ip = it[1]
	//        }
	//    } else {
	//        aclDto.ip = it
	//    }
	//    aclDto
	//}.filter {
	//    if (it.ip!!.isEmpty()) {
	//        return@filter false
	//    }
	//    ipRegex.matches(it.ip!!)
	//}.distinctBy { it.ip }
	//
	////添加黑白名单
	//ChannelAclDao.add(dto.id!!, aclDtoList)
	//
	//GlobalScope.launch(Dispatchers.IO) {
	//}
	//}

	//关闭正在通信的UDP连接
	//UDPBridgeManager.closeByChannel(channelId)

	//关闭代理监听
	proxy.CloseByChannel(channel.Id)
	client := ClientDao.SelectOne(channel.ClientId)

	//重新开启监听该客户端
	proxy.AcceptClient(client)
	return nil
}

// 表单验证
func validate(form form.ChannelEditForm) error {
	if len(form.Name) == 0 {
		return &controller.BusinessException{
			Message: "请填写隧道名",
		}
	}
	if len(form.Name) > 32 {
		return &controller.BusinessException{
			Message: "隧道名长度不能超过32个字符",
		}
	}
	//if len(form.ServerPort) == 0 {
	//	return &controller.BusinessException{
	//		Message: "服务端口必须设置",
	//	}
	//}
	//port, err := strconv.ParseInt(form.ServerPort, 10, 64)
	//if err != nil {
	//	return &controller.BusinessException{
	//		Message: "服务端口必须是一个数字",
	//	}
	//}
	if form.ServerPort < 0 || form.ServerPort > 65535 {
		return &controller.BusinessException{
			Message: "服务端口必须在0到65535之间",
		}
	}
	portChannel := ChannelDao.SelectByPort(form.ServerPort)
	if form.Id == 0 { //创建时
		if portChannel != nil {
			return &controller.BusinessException{
				Message: fmt.Sprintf("端口:%d已经被其他隧道占用", form.ServerPort),
			}
		}
	} else {
		if portChannel != nil && portChannel.Id != form.Id {
			return &controller.BusinessException{
				Message: fmt.Sprintf("端口:%d已经被其他隧道占用", form.ServerPort),
			}
		}
	}
	return nil
}
