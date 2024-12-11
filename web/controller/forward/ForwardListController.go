package forward

import (
	"DairoNPS/dao/DateDataSizeDao"
	"DairoNPS/dao/ForwardDao"
	"DairoNPS/extension/Number"
	"DairoNPS/forward"
	"DairoNPS/web/controller/forward/form"
)

// get:/forward_list
// templates:forward_list.html
func InitList() {
}

// 获取所有转发列表
// post:/forward_list/get_list
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
// post:/forward_list/delete
func Delete(id int) {

	//关闭隧道之后再打开
	forward.CloseAccept(id)
	DateDataSizeDao.DeleteByForward(id)
	ForwardDao.Delete(id)
}

// 修改可用状态
// post:/forward_list/set_state
func SetState(id int) {
	forwardDto := ForwardDao.SelectOne(id)
	if forwardDto.EnableState == 0 {
		ForwardDao.SetEnableState(id, 1)
		forward.Accept(ForwardDao.SelectOne(id))
	} else {
		ForwardDao.SetEnableState(id, 0)
		forward.CloseAccept(id)
	}
}
