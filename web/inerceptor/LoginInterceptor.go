package inerceptor

import (
	"DairoNPS/web/login_state"
	"net/http"
)

// LoginValidate 管理员登录验证
// interceptor:before
// include:/**
// exclude:/login**
func LoginValidate(writer http.ResponseWriter, request *http.Request) bool {
	if login_state.IsLogin(request) { //如果已经登录
		return true
	}
	if request.Method == "GET" {
		http.Redirect(writer, request, "/login", http.StatusFound)
		return false
	}
	writer.WriteHeader(http.StatusInternalServerError) // 设置状态码
	writer.Write([]byte(`{"Code":5,"Message":"未登录。"}`))
	return false
}
