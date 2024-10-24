package login

import (
	"DairoNPS/constant/NPSConstant"
	"DairoNPS/web"
	"DairoNPS/web/controller"
	"DairoNPS/web/controller/login/form"
	"net/http"
)

// 记录密码错误次数
var loginErrorCount = 0

// 初始化
func init() {
	http.HandleFunc("/login/do_login", web.ApiHandler(doLogin))
	http.HandleFunc("/login/login_out", web.ApiHandler(logout))
}

// 登录API
func doLogin(inForm form.LoginForm) any {
	if loginErrorCount > 10 {
		return &controller.BusinessException{
			Message: "用户名或密码错误次数超过限制，请重启程序再试。",
		}
	}
	err := validate(inForm)
	if err != nil {
		return err
	}
	if inForm.Name != NPSConstant.LoginName || inForm.Pwd != NPSConstant.LoginPwd {
		loginErrorCount++
		return &controller.BusinessException{
			Message: "用户名或密码错误",
		}
	}
	loginErrorCount = 0
	NPSConstant.IsLogin = true
	return nil
}

// 退出登录
func logout() {
	NPSConstant.IsLogin = false
}

// 表单验证
func validate(inForm form.LoginForm) error {
	if len(inForm.Name) == 0 {
		return &controller.BusinessException{
			Message: "用户名不能为空",
		}
	}
	if len(inForm.Pwd) == 0 {
		return &controller.BusinessException{
			Message: "密码不能为空",
		}
	}
	return nil
}
