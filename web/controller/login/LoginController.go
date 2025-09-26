package login

				import (
					"DairoNPS/DebugTimer"
				)

import (
	"DairoNPS/constant/NPSConstant"
	"DairoNPS/extension/String"
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

// get:/login
// templates:login.html
func Login() {
DebugTimer.Add584()
}

// DoLogin 登录API
// post:/login/do_login
func DoLogin(writer http.ResponseWriter, inForm form.LoginForm) any {
DebugTimer.Add585()
	if loginErrorCount > 10 {
DebugTimer.Add586()
		return &controller.BusinessException{
			Message: "用户名或密码错误次数超过限制，请重启程序再试。",
		}
	}
	err := validate(inForm)
	if err != nil {
DebugTimer.Add587()
		return err
	}
	if inForm.Name != NPSConstant.LoginName || inForm.Pwd != NPSConstant.LoginPwd {
DebugTimer.Add588()
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

// Logout 退出登录
// post:/login/login_out
func Logout() {
DebugTimer.Add589()
	login_state.LoginOut()
}

// 表单验证
func validate(inForm form.LoginForm) error {
DebugTimer.Add590()
	if len(inForm.Name) == 0 {
DebugTimer.Add591()
		return &controller.BusinessException{
			Message: "用户名不能为空",
		}
	}
	if len(inForm.Pwd) == 0 {
DebugTimer.Add592()
		return &controller.BusinessException{
			Message: "密码不能为空",
		}
	}
	return nil
}

// Logout 退出登录
// post:/login/login_out/test
func LogoutTest() string {
DebugTimer.Add593()
	return "123"
}
