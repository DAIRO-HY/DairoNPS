package forward

import (
	"DairoNPS/dao/ForwardDao"
	"DairoNPS/extension/Number"
	"DairoNPS/forward"
	"DairoNPS/web"
	"DairoNPS/web/controller/forward/form"
	"net/http"
)

// 初始化
func init() {
	http.HandleFunc("/forward_list/get_list", web.ApiHandler(GetList))
	http.HandleFunc("/forward_list/delete", web.ApiHandler(Delete))
	http.HandleFunc("/forward_list/set_state", web.ApiHandler(SetState))

}

// 获取所有转发列表
func GetList() []form.ForwardListOutForm {
	list := ForwardDao.SelectAll()
	outFormList := make([]form.ForwardListOutForm, len(list))
	for i, it := range list {
		outFormList[i] = form.ForwardListOutForm{
			Id:          it.Id,
			Name:        it.Name,
			Port:        it.Port,
			TargetPort:  it.TargetPort,
			EnableState: it.EnableState,
			InData:      Number.ToDataSize(it.InData),
			OutData:     Number.ToDataSize(it.OutData),
			Error:       it.Error,
		}
	}
	return outFormList
}

// 通过id删除一个转发
func Delete(inForm form.ForwardDeleteInputForm) {

	//关闭隧道之后再打开
	forward.CloseAccept(inForm.Id)
	ForwardDao.Delete(inForm.Id)
}

// 修改可用状态
func SetState(inForm form.ForwardSetStateInputForm) {
	forwardDto := ForwardDao.SelectOne(inForm.Id)
	if forwardDto.EnableState == 0 {
		ForwardDao.SetEnableState(inForm.Id, 1)
		forward.Accept(ForwardDao.SelectOne(inForm.Id))
	} else {
		ForwardDao.SetEnableState(inForm.Id, 0)
		forward.CloseAccept(inForm.Id)
	}
}
