package inerceptor

				import (
					"DairoNPS/DebugTimer"
				)

import (
	"DairoNPS/web/login_state"
	"net/http"
)

// LoginValidate 管理员登录验证
// interceptor:before
// include:/**
// exclude:/login**
func LoginValidate(writer http.ResponseWriter, request *http.Request) bool {
DebugTimer.Add601()
	if login_state.IsLogin(request) { //如果已经登录
DebugTimer.Add602()
		return true
	}
	if request.Method == "GET" {
DebugTimer.Add603()
		http.Redirect(writer, request, "/login", http.StatusFound)
		return false
	}
	writer.WriteHeader(http.StatusInternalServerError) // 设置状态码
	writer.Write([]byte(`{"Code":5,"Message":"未登录。"}`))
	return false
}
