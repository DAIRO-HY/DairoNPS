package login

import (
	"DairoNPS/constant/NPSConstant"
	"DairoNPS/extension/String"
	"DairoNPS/web"
	"DairoNPS/web/controller"
	"DairoNPS/web/controller/login/form"
	"DairoNPS/web/login_state"
	"math/rand/v2"
	"net/http"
	"strconv"
	"time"
)

// 记录密码错误次数
var loginErrorCount = 0

// 初始化
func init() {
	http.HandleFunc("/login/do_login", web.ApiHandler(doLogin))
	http.HandleFunc("/login/login_out", web.ApiHandler(logout))
}

// 登录API
func doLogin(writer http.ResponseWriter, request *http.Request, inForm form.LoginForm) any {
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

	timeRand := time.Now().UnixMilli() + int64(rand.IntN(900000)+100000)
	timeRandStr := strconv.FormatInt(timeRand, 10)
	token := String.ToMd5(timeRandStr)
	tokenCookie := &http.Cookie{
		Name:    login_state.COOKIE_TOKEN,
		Value:   token,
		Path:    "/",
		Expires: time.Now().AddDate(100, 0, 0), //100年以后过期
		MaxAge:  100 * 365 * 24 * 60 * 60,
		//HttpOnly: true,
	}
	http.SetCookie(writer, tokenCookie)
	
	login_state.Login(token)
	return nil
}

// 退出登录
func logout() {
	login_state.LoginOut()
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
