package forward

				import (
					"DairoNPS/DebugTimer"
				)

import (
	"DairoNPS/dao/ChannelDao"
	"DairoNPS/dao/ForwardDao"
	"DairoNPS/dao/dto"
	"DairoNPS/extension/Bool"
	"DairoNPS/extension/Date"
	"DairoNPS/extension/Number"
	forwardTcp "DairoNPS/forward"
	"DairoNPS/web/controller"
	"DairoNPS/web/controller/forward/form"
	"fmt"
)

// get:/forward_list/forward_edit
// templates:forward_edit.html
func InitEdit() {
DebugTimer.Add553()
}

// 获取转发信息
// post:/forward_list/forward_edit/info
func Info(id int) form.ForwardEditForm {
DebugTimer.Add554()
	if id != 0 {
DebugTimer.Add555()
		forwardDto := ForwardDao.SelectOne(id)
		return form.ForwardEditForm{
			Id:          forwardDto.Id,
			Name:        forwardDto.Name,
			Remark:      forwardDto.Remark,
			Port:        forwardDto.Port,
			TargetPort:  forwardDto.TargetPort,
			Date:        Date.FormatByTimespan(forwardDto.Date),
			EnableState: Bool.Is(forwardDto.EnableState == 0, "关闭", "开启"),
			InData:      Number.ToDataSize(forwardDto.InData),
			OutData:     Number.ToDataSize(forwardDto.OutData),
			AclState:    forwardDto.AclState,
		}
	}
	return form.ForwardEditForm{}
}

// 提交表单API
// post:/forward_list/forward_edit/edit
func Edit(inForm form.ForwardEditForm) any {
DebugTimer.Add556()
	err := validate(inForm)
	if err != nil {
DebugTimer.Add557()
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
DebugTimer.Add558()
		ForwardDao.Add(forwardDto)
	} else { //更新时
		ForwardDao.Update(forwardDto)
	}

	newDto := ForwardDao.SelectOne(forwardDto.Id)
	forwardTcp.CloseAccept(newDto.Id)
	if newDto.EnableState == 1 {
DebugTimer.Add559()
		forwardTcp.Accept(newDto)
	}
	return nil
}

/**
 * 表单验证
 */
func validate(inForm form.ForwardEditForm) error {
DebugTimer.Add560()
	if len(inForm.Name) == 0 {
DebugTimer.Add561()
		return &controller.BusinessException{
			Message: "请填写名称",
		}
	}
	if len(inForm.Name) > 32 {
DebugTimer.Add562()
		return &controller.BusinessException{
			Message: "名称长度不能超过32位",
		}
	}

	if inForm.Port < 0 || inForm.Port > 65535 {
DebugTimer.Add563()
		return &controller.BusinessException{
			Message: "端口必须在0到65535之间",
		}
	}
	portDto := ForwardDao.SelectByPort(inForm.Port)
	if inForm.Id == 0 { //创建时
DebugTimer.Add564()
		if portDto != nil {
DebugTimer.Add565()
			return &controller.BusinessException{
				Message: fmt.Sprintf("端口:%d已经被%s占用", portDto.Port, portDto.Name),
			}
		}
	} else {
		if portDto != nil && portDto.Id != inForm.Id {
DebugTimer.Add566()
			return &controller.BusinessException{
				Message: fmt.Sprintf("端口:%d已经被%s占用", portDto.Port, portDto.Name),
			}
		}
	}
	portChannel := ChannelDao.SelectByPort(inForm.Port)
	if portChannel != nil {
DebugTimer.Add567()
		return &controller.BusinessException{
			Message: fmt.Sprintf("端口:%d已被隧道:%s 占用", portChannel.ServerPort, portChannel.Name),
		}
	}
	return nil
}
