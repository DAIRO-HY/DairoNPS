package main

import (
	"DairoNPS/constant/NPSConstant"
	"DairoNPS/forward"
	"DairoNPS/nps/nps_client/tcp_client"
	"DairoNPS/nps/nps_client/udp_client"
	"DairoNPS/nps/nps_pool/tcp_pool"
	"DairoNPS/nps/nps_pool/udp_pool"
	"DairoNPS/util/LogUtil"
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
const VERSION = "1.0.3"

func init() {

	// 初始化共享接口
	tcp_pool.Csmi = &tcp_client.ClientSessionManager{}
	udp_pool.Csmi = &tcp_client.ClientSessionManager{}
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
	tcp_client.Accept()
	udp_client.Accept()

}

// 解析参数
func parseArgs() {
	fmt.Println("------------------------------------------------------------------------")
	for _, it := range os.Args {
		fmt.Println(it)
	}
	fmt.Println("------------------------------------------------------------------------")
	for _, it := range os.Args {
		paramArr := strings.Split(it, ":")
		switch paramArr[0] {
		case "-login-name":
			NPSConstant.LoginName = paramArr[1]
		case "-login-pwd":
			if len(paramArr[1]) > 0 {
				NPSConstant.LoginPwd = paramArr[1]
			}
		case "-web-port":
			NPSConstant.WebPort = paramArr[1]
		case "-tcp-port":
			NPSConstant.TcpPort = paramArr[1]
		case "-log-type": //日志输出方式
			switch paramArr[1] {
			case "0":
				LogUtil.LogOutType = LogUtil.LOG_OUT_TYPE_NO
			case "1":
				LogUtil.LogOutType = LogUtil.LOG_OUT_TYPE_CONSOLE
			case "2":
				LogUtil.LogOutType = LogUtil.LOG_OUT_TYPE_FILE
			}
		case "-log-level": //日志输出级别
			levels := strings.Split(paramArr[1], ",")
			for _, level := range levels {
				LogUtil.LogLevel[level] = true
			}
		}
	}
	fmt.Printf("程序启动成功，管理员：%s 密码：%s\n", NPSConstant.LoginName, NPSConstant.LoginPwd)
}
