package channel

import (
	"DairoNPS/dao/ChannelDao"
	"DairoNPS/dao/ClientDao"
	"DairoNPS/dao/dto"
	"DairoNPS/extension/Bool"
	"DairoNPS/extension/Date"
	"DairoNPS/web"
	"DairoNPS/web/controller"
	"DairoNPS/web/controller/channel/form"
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
		outForm = form.ChannelEditForm{}
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
			InDataTotal:   channelDto.InDataTotal,
			OutDataTotal:  channelDto.OutDataTotal,
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
	//
	//    //关闭隧道之后再打开
	//    CLServer.closeByChannel(dto.id!!)
	//    CLServer.start(dto.id!!)
	//}
	//} catch (e: Exception) {
	//    e.message ?: throw e
	//    if (e.message!!.contains("UNIQUE constraint failed: channel.server_port")) {
	//        throw BusinessException("服务端口已被占用")
	//    }
	//    throw e
	//}
	return nil
}

// 表单验证
func validate(editForm form.ChannelEditForm) error {
	if len(editForm.Name) == 0 {
		return &controller.BusinessException{
			Message: "请填写隧道名",
		}
	}
	if len(editForm.Name) > 32 {
		return &controller.BusinessException{
			Message: "隧道名长度不能超过32个字符",
		}
	}
	return nil
}

//}
