package forward

import (
	"DairoNPS/dao/ChannelDao"
	"DairoNPS/dao/ForwardDao"
	"DairoNPS/dao/dto"
	"DairoNPS/extension/Bool"
	"DairoNPS/extension/Date"
	"DairoNPS/extension/Number"
	forwardTcp "DairoNPS/forward"
	"DairoNPS/web"
	"DairoNPS/web/controller"
	"DairoNPS/web/controller/forward/form"
	"fmt"
	"net/http"
)

// 初始化
func init() {
	http.HandleFunc("/forward_list/forward_edit/info", web.ApiHandler(Info))
	http.HandleFunc("/forward_list/forward_edit/edit", web.ApiHandler(Edit))

}

func Info(inForm form.ForwardInfoInputForm) form.ForwardEditForm {
	if inForm.Id != 0 {
		forwardDto := ForwardDao.SelectOne(inForm.Id)
		return form.ForwardEditForm{
			Id:           forwardDto.Id,
			Name:         forwardDto.Name,
			Remark:       forwardDto.Remark,
			Port:         forwardDto.Port,
			TargetPort:   forwardDto.TargetPort,
			CreateDate:   Date.FormatByTimespan(forwardDto.Date),
			EnableState:  Bool.Is(forwardDto.EnableState == 0, "关闭", "开启"),
			InDataTotal:  Number.ToDataSize(forwardDto.InData),
			OutDataTotal: Number.ToDataSize(forwardDto.OutData),
			AclState:     forwardDto.AclState,
		}
	}
	return form.ForwardEditForm{}
}

// 提交表单API
func Edit(inForm form.ForwardEditForm) any {
	err := validate(inForm)
	if err != nil {
		return err
	}
	forwardDto := &dto.ForwardDto{
		Id:         inForm.Id,
		Name:       inForm.Name,
		Port:       inForm.Port,
		TargetPort: inForm.TargetPort,
		AclState:   inForm.AclState,
		Remark:     inForm.Remark,
	}
	if inForm.Id == 0 {
		ForwardDao.Add(forwardDto)
	} else { //更新时
		ForwardDao.Update(forwardDto)
	}

	newDto := ForwardDao.SelectOne(forwardDto.Id)
	forwardTcp.CloseAccept(newDto.Id)
	if newDto.EnableState == 1 {
		forwardTcp.Accept(newDto)
	}
	return nil

	//	//判断是否IP地址的正则表达式
	//	val ipRegex = """^(([01]?\d{1,2}|2[0-4]\d|25[0-5])\.){3}([01]?\d{1,2}|2[0-4]\d|25[0-5])$""".toRegex()
	//	val aclDtoList = form.aclIp.map{
	//	val aclDto = ForwardAclDto()
	//	if (it.contains(":")){
	//	it.split(":").also{
	//	aclDto.remark = it[0]
	//	aclDto.ip = it[1]
	//}
	//} else{
	//	aclDto.ip = it
	//}
	//	aclDto
	//}.filter{
	//	if (it.ip!!.isEmpty()){
	//	return @filter false
	//}
	//	ipRegex.matches(it.ip!!)
	//}.distinctBy{it.ip}
	//
	//	//添加黑白名单
	//	ForwardAclDao.add(dto.id!!, aclDtoList)
	//
	//	GlobalScope.launch(Dispatchers.IO){
	//
	//	//关闭隧道之后再打开
	//	//                CLServer.closeByChannel(dto.id!!)
	//	//                CLServer.start(dto.id!!)
	//
	//	ForwardServer.close(dto.id!!)
	//	ForwardServer.start(dto.id!!)
	//}
	//}
	//catch(e: Exception) {
	//	e.message ?: throw
	//	e
	//	if (e.message!!.contains("UNIQUE constraint failed: forward.server_port")) {
	//		throw
	//		BusinessException("服务端口已被占用")
	//	}
	//	throw
	//	e
	//}
}

/**
 * 表单验证
 */
func validate(inForm form.ForwardEditForm) error {
	if len(inForm.Name) == 0 {
		return &controller.BusinessException{
			Message: "请填写名称",
		}
	}
	if len(inForm.Name) > 32 {
		return &controller.BusinessException{
			Message: "名称长度不能超过32位",
		}
	}

	if inForm.Port < 0 || inForm.Port > 65535 {
		return &controller.BusinessException{
			Message: "端口必须在0到65535之间",
		}
	}
	portDto := ForwardDao.SelectByPort(inForm.Port)
	if inForm.Id == 0 { //创建时
		if portDto != nil {
			return &controller.BusinessException{
				Message: fmt.Sprintf("端口:%d已经被%s占用", portDto.Port, portDto.Name),
			}
		}
	} else {
		if portDto != nil && portDto.Id != inForm.Id {
			return &controller.BusinessException{
				Message: fmt.Sprintf("端口:%d已经被%s占用", portDto.Port, portDto.Name),
			}
		}
	}
	portChannel := ChannelDao.SelectByPort(inForm.Port)
	if portChannel != nil {
		return &controller.BusinessException{
			Message: fmt.Sprintf("端口:%d已被隧道:%s 占用", portChannel.ServerPort, portChannel.Name),
		}
	}
	return nil
}