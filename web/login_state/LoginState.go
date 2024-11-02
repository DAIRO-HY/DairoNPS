package login_state

import (
	"DairoNPS/constant/NPSConstant"
	"net/http"
	"time"
)

type LoginState struct {

	//客户端IP
	ip string

	//登录token
	token string

	//登录时间
	loginDate int64

	//最后在线时间
	onlineDate int64
}

const COOKIE_TOKEN = "session_token"

// 登录状态
var loginState *LoginState

// 是否登录验证
func IsLogin(request *http.Request) bool {
	if NPSConstant.IsDev { //开发模式不需要登录
		return true
	}
	if loginState == nil {
		return false
	}
	cookie, _ := request.Cookie(COOKIE_TOKEN)
	if cookie == nil {
		return false
	}
	if loginState.token == cookie.Value {

		//更新最后在线时间
		loginState.onlineDate = time.Now().UnixMilli()
		return true
	}
	return false
}

// 设置登录状态
func Login(token string) {
	loginState = &LoginState{

		//客户端IP
		//ip string

		//登录token
		token: token,

		//登录时间
		loginDate: time.Now().UnixMilli(),

		//最后在线时间
		onlineDate: time.Now().UnixMilli(),
	}
}

// 退出登录
func LoginOut() {
	loginState = nil
}
