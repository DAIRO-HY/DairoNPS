package client

import (
	"DairoNPS/dao/ClientDao"
	"DairoNPS/dao/dto"
	"DairoNPS/extension/Bool"
	"DairoNPS/extension/Date"
	"DairoNPS/extension/Number"
	"DairoNPS/nps/nps_client"
	"DairoNPS/web"
	"DairoNPS/web/controller"
	"DairoNPS/web/controller/client/form"
	"net/http"
)

// 初始化
func init() {
	http.HandleFunc("/client_list/client_edit/info", web.ApiHandler(Info))
	http.HandleFunc("/client_list/client_edit/edit", web.ApiHandler(Edit))
}

func Info(inForm form.ClientEditForm) form.ClientEditForm {
	if inForm.Id != 0 {
		clientDto := ClientDao.SelectOne(inForm.Id)
		return form.ClientEditForm{
			Id:            clientDto.Id,
			Name:          clientDto.Name,
			Key:           clientDto.Key,
			Remark:        clientDto.Remark,
			Ip:            clientDto.Ip,
			Version:       clientDto.Version,
			Date:          Date.FormatByTimespan(clientDto.Date),
			EnableState:   Bool.Is(clientDto.EnableState == 0, "关闭", "开启"),
			InData:        Number.ToDataSize(clientDto.InData),
			OutData:       Number.ToDataSize(clientDto.OutData),
			OnlineState:   Bool.Is(nps_client.IsOnline(clientDto.Id), "在线", "离线"),
			LastLoginDate: Date.FormatByTimespan(clientDto.LastLoginDate),
		}
	}
	return form.ClientEditForm{}
}

// 提交表单API
func Edit(form form.ClientEditForm) any {
	err := validate(form)
	if err != nil {
		return err
	}
	var clientDto *dto.ClientDto
	if form.Id != 0 {
		clientDto = ClientDao.SelectOne(form.Id)
	} else {
		clientDto = &dto.ClientDto{}
	}

	//名称
	clientDto.Name = form.Name

	//版本号
	clientDto.Version = form.Version

	//连接认证秘钥
	clientDto.Key = form.Key

	//一些备注信息,错误信息等
	clientDto.Remark = form.Remark
	if form.Id == 0 {
		ClientDao.Add(clientDto)
	} else {
		ClientDao.Update(clientDto)
	}
	nps_client.Shutdown(form.Id)
	return nil
}

/**
 * 表单验证
 */
func validate(form form.ClientEditForm) error {
	if len(form.Name) == 0 {
		return &controller.BusinessException{
			Message: "请填写名称",
		}
	}
	if len(form.Name) > 32 {
		return &controller.BusinessException{
			Message: "名称长度不能超过32位",
		}
	}
	if len(form.Key) == 0 {
		return &controller.BusinessException{
			Message: "请填写认证秘钥",
		}
	}
	if len(form.Key) > 32 {
		return &controller.BusinessException{
			Message: "认证秘钥长度不得超过32位",
		}
	}
	keyClient := ClientDao.SelectByKey(form.Key)
	if form.Id == 0 { //添加数据时
		if keyClient != nil {
			return &controller.BusinessException{
				Message: "该秘钥已被其他客户端使用，请换一个秘钥。",
			}
		}
	} else { //更新时
		if keyClient != nil && keyClient.Id != form.Id {
			return &controller.BusinessException{
				Message: "该秘钥已被其他客户端使用，请换一个秘钥。",
			}
		}
	}
	return nil
}
