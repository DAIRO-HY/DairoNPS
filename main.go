package main

import (
	"DairoNPS/constant/NPSConstant"
	"DairoNPS/forward"
	"DairoNPS/nps/nps_client"
	"DairoNPS/nps/nps_pool"
	"DairoNPS/web"
	"fmt"
	"os"
	"strings"

	//初始化Controller
	_ "DairoNPS/web/controller/bridge_list"
	_ "DairoNPS/web/controller/channel"
	_ "DairoNPS/web/controller/client"
	_ "DairoNPS/web/controller/common"
	_ "DairoNPS/web/controller/data_size_log"
	_ "DairoNPS/web/controller/forward"
	_ "DairoNPS/web/controller/index"
	_ "DairoNPS/web/controller/login"
	_ "DairoNPS/web/controller/speed_chart"
)

// 版本号
const VERSION = "1.0.0"

func init() {

	// 初始化共享接口
	nps_pool.Csmi = &nps_client.ClientSessionManager{}
}

var list = make([]int, 0)

func main() {

	// 解析参数
	parseArgs()

	//启动web管理
	go web.Start()

	//启动端口转发
	go forward.StartAcceptAll()

	//启动客户端监听
	nps_client.Accept()

}

// 解析参数
func parseArgs() {
	for _, it := range os.Args {
		paramArr := strings.Split(it, ":")
		switch paramArr[0] {
		case "-login-name":
			NPSConstant.LoginName = paramArr[1]
		case "-login-pwd":
			NPSConstant.LoginPwd = paramArr[1]
		case "-web-port":
			NPSConstant.WebPort = paramArr[1]
		case "-tcp-port":
			NPSConstant.TcpPort = paramArr[1]
		}
	}
	fmt.Printf("程序启动成功，管理员：%s 密码：%s\n", NPSConstant.LoginName, NPSConstant.LoginPwd)
}
