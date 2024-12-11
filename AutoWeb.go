/**
 * 代码为自动生成，请勿手动修改
 */
package main

import (
	webcontrollerbridgelist "DairoNPS/web/controller/bridge_list"
	webcontrollerbridgelistform "DairoNPS/web/controller/bridge_list/form"
	webcontrollerchannel "DairoNPS/web/controller/channel"
	webcontrollerchannelform "DairoNPS/web/controller/channel/form"
	webcontrollerclient "DairoNPS/web/controller/client"
	webcontrollerclientform "DairoNPS/web/controller/client/form"
	webcontrollercommon "DairoNPS/web/controller/common"
	webcontrollerdatasizelog "DairoNPS/web/controller/data_size_log"
	webcontrollerdatasizelogform "DairoNPS/web/controller/data_size_log/form"
	webcontrollerforward "DairoNPS/web/controller/forward"
	webcontrollerforwardform "DairoNPS/web/controller/forward/form"
	webcontrollerindex "DairoNPS/web/controller/index"
	webcontrollerlogin "DairoNPS/web/controller/login"
	webcontrollerloginform "DairoNPS/web/controller/login/form"
	webcontrollerspeedchart "DairoNPS/web/controller/speed_chart"
	webinerceptor "DairoNPS/web/inerceptor"

	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

//go:embed resources/static/*
var staticFiles embed.FS

//go:embed resources/templates/*
var templatesFiles embed.FS

// 开启web服务
func startWebServer(port int) {

	// 将嵌入的资源限制到 "/resources/static" 子目录
	staticFS, staticErr := fs.Sub(staticFiles, "resources/static")
	if staticErr != nil {
		panic(staticErr)
	}

	// 使用 http.FileServer 提供文件服务
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	http.HandleFunc("/bridge_list", func(writer http.ResponseWriter, request *http.Request) {
		if !webinerceptor.LoginValidate(writer, request) {
			return
		}
		var body any = nil
		webcontrollerbridgelist.Init()
		templates := append([]string{"resources/templates/bridge_list.html"}, COMMON_TEMPLATES...)
		writeToTemplate(writer, templates, body)
	})
	http.HandleFunc("/bridge_list/load_data", func(writer http.ResponseWriter, request *http.Request) {
		if !webinerceptor.LoginValidate(writer, request) {
			return
		}
		paramMap := makeParamMap(request)
		search := getForm[webcontrollerbridgelistform.BridgeInForm](paramMap)
		var body any = nil
		body = webcontrollerbridgelist.LoadData(search)
		writeToResponse(writer, body)
	})
	http.HandleFunc("/channel_list/channel_edit", func(writer http.ResponseWriter, request *http.Request) {
		if !webinerceptor.LoginValidate(writer, request) {
			return
		}
		var body any = nil
		webcontrollerchannel.InitEdit()
		templates := append([]string{"resources/templates/channel_edit.html"}, COMMON_TEMPLATES...)
		writeToTemplate(writer, templates, body)
	})
	http.HandleFunc("/channel_list/channel_edit/info", func(writer http.ResponseWriter, request *http.Request) {
		if !webinerceptor.LoginValidate(writer, request) {
			return
		}
		paramMap := makeParamMap(request)
		ClientId := getInt(paramMap, "ClientId")
		Id := getInt(paramMap, "Id")
		var body any = nil
		body = webcontrollerchannel.Info(ClientId, Id)
		writeToResponse(writer, body)
	})
	http.HandleFunc("/channel_list/channel_edit/edit", func(writer http.ResponseWriter, request *http.Request) {
		if !webinerceptor.LoginValidate(writer, request) {
			return
		}
		paramMap := makeParamMap(request)
		form := getForm[webcontrollerchannelform.ChannelEditForm](paramMap)
		var body any = nil
		body = webcontrollerchannel.Edit(form)
		writeToResponse(writer, body)
	})
	http.HandleFunc("/channel_list", func(writer http.ResponseWriter, request *http.Request) {
		if !webinerceptor.LoginValidate(writer, request) {
			return
		}
		var body any = nil
		webcontrollerchannel.InitList()
		templates := append([]string{"resources/templates/channel_list.html"}, COMMON_TEMPLATES...)
		writeToTemplate(writer, templates, body)
	})
	http.HandleFunc("/channel_list/list", func(writer http.ResponseWriter, request *http.Request) {
		if !webinerceptor.LoginValidate(writer, request) {
			return
		}
		paramMap := makeParamMap(request)
		clientId := getInt(paramMap, "clientId")
		var body any = nil
		body = webcontrollerchannel.List(clientId)
		writeToResponse(writer, body)
	})
	http.HandleFunc("/channel_list/delete", func(writer http.ResponseWriter, request *http.Request) {
		if !webinerceptor.LoginValidate(writer, request) {
			return
		}
		paramMap := makeParamMap(request)
		id := getInt(paramMap, "id")
		var body any = nil
		webcontrollerchannel.Delete(id)
		writeToResponse(writer, body)
	})
	http.HandleFunc("/channel_list/set_state", func(writer http.ResponseWriter, request *http.Request) {
		if !webinerceptor.LoginValidate(writer, request) {
			return
		}
		paramMap := makeParamMap(request)
		id := getInt(paramMap, "id")
		var body any = nil
		webcontrollerchannel.SetState(id)
		writeToResponse(writer, body)
	})
	http.HandleFunc("/client_list/client_edit", func(writer http.ResponseWriter, request *http.Request) {
		if !webinerceptor.LoginValidate(writer, request) {
			return
		}
		var body any = nil
		webcontrollerclient.InitEdit()
		templates := append([]string{"resources/templates/client_edit.html"}, COMMON_TEMPLATES...)
		writeToTemplate(writer, templates, body)
	})
	http.HandleFunc("/client_list/client_edit/info", func(writer http.ResponseWriter, request *http.Request) {
		if !webinerceptor.LoginValidate(writer, request) {
			return
		}
		paramMap := makeParamMap(request)
		id := getInt(paramMap, "id")
		var body any = nil
		body = webcontrollerclient.Info(id)
		writeToResponse(writer, body)
	})
	http.HandleFunc("/client_list/client_edit/edit", func(writer http.ResponseWriter, request *http.Request) {
		if !webinerceptor.LoginValidate(writer, request) {
			return
		}
		paramMap := makeParamMap(request)
		form := getForm[webcontrollerclientform.ClientEditForm](paramMap)
		var body any = nil
		body = webcontrollerclient.Edit(form)
		writeToResponse(writer, body)
	})
	http.HandleFunc("/client_list", func(writer http.ResponseWriter, request *http.Request) {
		if !webinerceptor.LoginValidate(writer, request) {
			return
		}
		var body any = nil
		webcontrollerclient.InitList()
		templates := append([]string{"resources/templates/client_list.html"}, COMMON_TEMPLATES...)
		writeToTemplate(writer, templates, body)
	})
	http.HandleFunc("/client_list/init", func(writer http.ResponseWriter, request *http.Request) {
		if !webinerceptor.LoginValidate(writer, request) {
			return
		}
		var body any = nil
		body = webcontrollerclient.List()
		writeToResponse(writer, body)
	})
	http.HandleFunc("/client_list/delete", func(writer http.ResponseWriter, request *http.Request) {
		if !webinerceptor.LoginValidate(writer, request) {
			return
		}
		paramMap := makeParamMap(request)
		id := getInt(paramMap, "id")
		var body any = nil
		webcontrollerclient.Delete(id)
		writeToResponse(writer, body)
	})
	http.HandleFunc("/client_list/set_state", func(writer http.ResponseWriter, request *http.Request) {
		if !webinerceptor.LoginValidate(writer, request) {
			return
		}
		paramMap := makeParamMap(request)
		id := getInt(paramMap, "id")
		var body any = nil
		webcontrollerclient.SetState(id)
		writeToResponse(writer, body)
	})
	http.HandleFunc("/common/dropdown", func(writer http.ResponseWriter, request *http.Request) {
		if !webinerceptor.LoginValidate(writer, request) {
			return
		}
		var body any = nil
		body = webcontrollercommon.Dropdown(request)
		writeToResponse(writer, body)
	})
	http.HandleFunc("/data_size/get_data_size", func(writer http.ResponseWriter, request *http.Request) {
		if !webinerceptor.LoginValidate(writer, request) {
			return
		}
		paramMap := makeParamMap(request)
		inForm := getForm[webcontrollerdatasizelogform.GetDataInForm](paramMap)
		var body any = nil
		body = webcontrollerdatasizelog.GetDataSize(inForm)
		writeToResponse(writer, body)
	})
	http.HandleFunc("/forward_list/forward_edit", func(writer http.ResponseWriter, request *http.Request) {
		if !webinerceptor.LoginValidate(writer, request) {
			return
		}
		var body any = nil
		webcontrollerforward.InitEdit()
		templates := append([]string{"resources/templates/forward_edit.html"}, COMMON_TEMPLATES...)
		writeToTemplate(writer, templates, body)
	})
	http.HandleFunc("/forward_list/forward_edit/info", func(writer http.ResponseWriter, request *http.Request) {
		if !webinerceptor.LoginValidate(writer, request) {
			return
		}
		paramMap := makeParamMap(request)
		id := getInt(paramMap, "id")
		var body any = nil
		body = webcontrollerforward.Info(id)
		writeToResponse(writer, body)
	})
	http.HandleFunc("/forward_list/forward_edit/edit", func(writer http.ResponseWriter, request *http.Request) {
		if !webinerceptor.LoginValidate(writer, request) {
			return
		}
		paramMap := makeParamMap(request)
		inForm := getForm[webcontrollerforwardform.ForwardEditForm](paramMap)
		var body any = nil
		body = webcontrollerforward.Edit(inForm)
		writeToResponse(writer, body)
	})
	http.HandleFunc("/forward_list", func(writer http.ResponseWriter, request *http.Request) {
		if !webinerceptor.LoginValidate(writer, request) {
			return
		}
		var body any = nil
		webcontrollerforward.InitList()
		templates := append([]string{"resources/templates/forward_list.html"}, COMMON_TEMPLATES...)
		writeToTemplate(writer, templates, body)
	})
	http.HandleFunc("/forward_list/get_list", func(writer http.ResponseWriter, request *http.Request) {
		if !webinerceptor.LoginValidate(writer, request) {
			return
		}
		var body any = nil
		body = webcontrollerforward.GetList()
		writeToResponse(writer, body)
	})
	http.HandleFunc("/forward_list/delete", func(writer http.ResponseWriter, request *http.Request) {
		if !webinerceptor.LoginValidate(writer, request) {
			return
		}
		paramMap := makeParamMap(request)
		id := getInt(paramMap, "id")
		var body any = nil
		webcontrollerforward.Delete(id)
		writeToResponse(writer, body)
	})
	http.HandleFunc("/forward_list/set_state", func(writer http.ResponseWriter, request *http.Request) {
		if !webinerceptor.LoginValidate(writer, request) {
			return
		}
		paramMap := makeParamMap(request)
		id := getInt(paramMap, "id")
		var body any = nil
		webcontrollerforward.SetState(id)
		writeToResponse(writer, body)
	})
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		if !webinerceptor.LoginValidate(writer, request) {
			return
		}
		var body any = nil
		webcontrollerindex.Home()
		templates := append([]string{"resources/templates/index.html"}, COMMON_TEMPLATES...)
		writeToTemplate(writer, templates, body)
	})
	http.HandleFunc("/index", func(writer http.ResponseWriter, request *http.Request) {
		if !webinerceptor.LoginValidate(writer, request) {
			return
		}
		var body any = nil
		webcontrollerindex.Init()
		templates := append([]string{"resources/templates/index.html"}, COMMON_TEMPLATES...)
		writeToTemplate(writer, templates, body)
	})
	http.HandleFunc("/index/data", func(writer http.ResponseWriter, request *http.Request) {
		if !webinerceptor.LoginValidate(writer, request) {
			return
		}
		var body any = nil
		webcontrollerindex.Data(writer, request)
		writeToResponse(writer, body)
	})
	http.HandleFunc("/index/gc", func(writer http.ResponseWriter, request *http.Request) {
		if !webinerceptor.LoginValidate(writer, request) {
			return
		}
		var body any = nil
		webcontrollerindex.Gc()
		writeToResponse(writer, body)
	})
	http.HandleFunc("/login", func(writer http.ResponseWriter, request *http.Request) {
		var body any = nil
		webcontrollerlogin.Login()
		templates := append([]string{"resources/templates/login.html"}, COMMON_TEMPLATES...)
		writeToTemplate(writer, templates, body)
	})
	http.HandleFunc("/login/do_login", func(writer http.ResponseWriter, request *http.Request) {
		paramMap := makeParamMap(request)
		inForm := getForm[webcontrollerloginform.LoginForm](paramMap)
		var body any = nil
		body = webcontrollerlogin.DoLogin(writer, inForm)
		writeToResponse(writer, body)
	})
	http.HandleFunc("/login/login_out", func(writer http.ResponseWriter, request *http.Request) {
		var body any = nil
		webcontrollerlogin.Logout()
		writeToResponse(writer, body)
	})
	http.HandleFunc("/login/login_out/test", func(writer http.ResponseWriter, request *http.Request) {
		var body any = nil
		body = webcontrollerlogin.LogoutTest()
		writeToResponse(writer, body)
	})
	http.HandleFunc("/ws/speed_chart", func(writer http.ResponseWriter, request *http.Request) {
		if !webinerceptor.LoginValidate(writer, request) {
			return
		}
		var body any = nil
		webcontrollerspeedchart.CurrentData(writer, request)
		writeToResponse(writer, body)
	})

	// 启动服务器
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}

// 生成参数Map
func makeParamMap(request *http.Request) map[string][]string {
	query := request.URL.Query()

	//解析post表单
	request.ParseForm()
	postParams := request.PostForm

	//将参数转换成Map
	paramMap := make(map[string][]string)
	for key, v := range query {
		paramMap[key] = v
	}
	for key, v := range postParams {
		paramMap[key] = v
	}
	return paramMap
}

// 获取表单实例
func getForm[T any](paramMap map[string][]string) T {

	// 创建结构体实例
	targetForm := new(T)
	reflectForm := reflect.ValueOf(targetForm).Elem()
	argType := reflect.TypeOf(*targetForm)

	// 遍历结构体字段
	for j := 0; j < argType.NumField(); j++ {
		field := argType.Field(j)
		fieldName := field.Name

		//得到参数值
		value := paramMap[fieldName]
		if value == nil {
			//将首字母小写再去获取参数
			lowerKey := strings.ToLower(fieldName[:1]) + fieldName[1:]
			value = paramMap[lowerKey]
		}
		if value == nil {
			continue
		}

		// 设置字段值（这里我们设置为示例值）
		switch field.Type.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

			// 设置整数字段
			intValue, _ := strconv.ParseInt(value[0], 10, 64)
			reflectForm.Field(j).SetInt(intValue)
		case reflect.Float32, reflect.Float64:
			floatValue, _ := strconv.ParseFloat(value[0], 64)
			reflectForm.Field(j).SetFloat(floatValue)
		case reflect.String:
			reflectForm.Field(j).SetString(value[0]) // 设置字符串字段
		}
	}
	return *targetForm
}

// 获取string类型的参数
func getString(paramMap map[string][]string, key string) string {
	value := paramMap[key]
	if value == nil {
		return ""
	}
	rValue := value[0]
	return rValue
}

// 获取int类型的参数
func getInt(paramMap map[string][]string, key string) int {
	value := paramMap[key]
	if value == nil {
		return 0
	}
	rValue, _ := strconv.Atoi(value[0])
	return rValue
}

// 获取int类型的参数
func getInt64(paramMap map[string][]string, key string) int64 {
	value := paramMap[key]
	if value == nil {
		return 0
	}
	rValue, _ := strconv.ParseInt(value[0], 10, 64)
	return rValue
}

// 获取float32类型的参数
func getFloat32(paramMap map[string][]string, key string) float32 {
	value := paramMap[key]
	if value == nil {
		return 0
	}
	rValue, _ := strconv.ParseFloat(value[0], 32)
	return float32(rValue)
}

// 获取float64类型的参数
func getFloat64(paramMap map[string][]string, key string) float64 {
	value := paramMap[key]
	if value == nil {
		return 0
	}
	rValue, _ := strconv.ParseFloat(value[0], 64)
	return rValue
}

// 返回结果
func writeToResponse(writer http.ResponseWriter, body any) {
	if body == nil {
		return
	}
	if body == "" {
		return
	}

	// 设置 Content-Type 头部信息
	writer.Header().Set("Content-Type", "text/plain;charset=UTF-8")

	switch returnBody := body.(type) {
	case string:
		writer.Write([]uint8(returnBody))
	case int:
		writer.Write([]uint8(strconv.Itoa(returnBody)))
	case int8:
		writer.Write([]uint8(strconv.Itoa(int(returnBody))))
	case int16:
		writer.Write([]uint8(strconv.Itoa(int(returnBody))))
	case int32:
		writer.Write([]uint8(strconv.Itoa(int(returnBody))))
	case int64:
		writer.Write([]uint8(strconv.FormatInt(returnBody, 10)))
	case error:
		// 设置 HTTP 状态码
		writer.WriteHeader(http.StatusInternalServerError) // 设置状态码
		jsonData, _ := json.Marshal(body)
		writer.Write(jsonData)
	default:
		jsonData, _ := json.Marshal(body)
		writer.Write(jsonData)
	}
}

// 写入html模板
func writeToTemplate(writer http.ResponseWriter, templates []string, data any) {

	// 解析嵌入的模板
	t, err := template.ParseFS(templatesFiles, templates...)
	if err != nil {
		http.Error(writer, fmt.Sprintf("Error loading template:%q", err), http.StatusInternalServerError)
		return
	}

	// 设置 Content-Type 头部信息
	writer.Header().Set("Content-Type", "text/html;charset=UTF-8")
	t.Execute(writer, data)
}
