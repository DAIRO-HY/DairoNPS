package web

import (
	"DairoNPS/constant/NPSConstant"
	"DairoNPS/util/LogUtil"
	"DairoNPS/web/login_state"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

//func renderTemplate(w http.ResponseWriter, tmpl123 string) {
//	//tmpl := template.New("base")
//	//tmpl, _ = tmpl.Parse(`{{define "T1"}}Hello, {{.Name}} from T1!{{end}}`)
//	//tmpl, _ = tmpl.Parse(`{{define "T2"}}Hi, {{.Name}} from T2!{{end}}`)
//	//tmpl, _ = tmpl.ParseFiles(`web/templates/index.html`)
//
//	tmpl, _ := template.ParseFiles(
//		//tmplPath,
//		"web/templates/index.html",
//		//"web/templates/include/content.html",
//		//"web/templates/include/css.html",
//		"web/templates/include/js.html",
//		//"web/templates/include/top-bar.html",
//	)
//
//	//data := struct {
//	//	Name string
//	//}{
//	//	Name: "Go Developer",
//	//}
//	tmpl.Execute(w, nil) // 选择 "T2" 模板来执行
//}

func renderTemplate(w http.ResponseWriter, tmpl string) {

	// 解析嵌入的模板
	t, err := template.ParseFS(templatesFiles,
		"templates/"+tmpl,
		"templates/include/head.html",
		"templates/include/top-bar.html",
		"templates/include/data_size_chart.html",
		"templates/include/speed_chart.html",
	)
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	// 设置 Content-Type 头部信息
	w.Header().Set("Content-Type", "text/html;charset=UTF-8")
	t.Execute(w, nil)
}

// 路由处理
func htmlHandler(writer http.ResponseWriter, request *http.Request) {
	paths := strings.Split(request.URL.Path, "/")
	htmlFile := paths[len(paths)-1]
	if len(htmlFile) == 0 {
		htmlFile = "index"
	}
	renderTemplate(writer, htmlFile+".html")
}

// api请求路由
func ApiHandler(controller any) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {

		// 设置 Content-Type 头部信息
		writer.Header().Set("Content-Type", "text/plain;charset=UTF-8")
		if !login_state.IsLogin(request) { //如果当前未登录
			path := request.URL.Path
			if !strings.HasPrefix(path, "/login/") { //非登录API，返回错误
				writer.WriteHeader(http.StatusInternalServerError) // 设置状态码
				writer.Write([]byte(`{"Code":5,"Message":"未登录。"}`))
				return
			}
		}

		// 获取函数的反射值
		fn := reflect.ValueOf(controller)

		// 获取函数的类型
		fnType := fn.Type()

		// 获取参数个数
		numArgs := fnType.NumIn()

		//controller需要的参数
		params := make([]reflect.Value, numArgs)

		// 遍历每个参数
		for i := 0; i < numArgs; i++ {
			argType := fnType.In(i)
			typeStr := argType.String()
			var value any
			if typeStr == "*http.Request" {
				value = request
			} else if typeStr == "http.ResponseWriter" {
				value = writer
			} else if strings.HasSuffix(typeStr, "Form") { //这是一个Form表单
				value = getForm(request, argType)
			}
			params[i] = reflect.ValueOf(value)
		}
		returnValues := fn.Call(params)
		if len(returnValues) == 0 {
			return
		}
		body := returnValues[0].Interface()
		if body == nil {
			return
		}

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
}

// 获取表单实例
func getForm(request *http.Request, argType reflect.Type) any {
	//如果该参数是一个Form表单

	// 创建结构体实例
	form := reflect.New(argType).Elem()
	query := request.URL.Query()

	//解析post表单
	request.ParseForm()
	postParams := request.PostForm

	//将参数转换成Map
	paramMap := make(map[string][]string)
	for key, v := range query {
		paramMap[strings.ToLower(key)] = v
	}
	for key, v := range postParams {
		paramMap[strings.ToLower(key)] = v
	}

	// 遍历结构体字段
	for j := 0; j < argType.NumField(); j++ {
		field := argType.Field(j)
		fieldName := field.Name

		//得到参数值
		value := paramMap[strings.ToLower(fieldName)]
		if value == nil {
			continue
		}

		// 设置字段值（这里我们设置为示例值）
		switch field.Type.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

			// 设置整数字段
			intValue, _ := strconv.ParseInt(value[0], 10, 64)
			form.Field(j).SetInt(intValue)
		case reflect.String:
			form.Field(j).SetString(value[0]) // 设置字符串字段
		}
	}
	return form.Interface()
}

// 通过这种方式将静态资源文件打包进二进制文件
//
//go:embed static/*
var staticFiles embed.FS

//go:embed templates/*
var templatesFiles embed.FS

func Start() {

	// 处理静态文件
	fs := http.FileServer(http.FS(staticFiles))
	http.Handle("/static/", fs)

	// 设置路由
	http.HandleFunc("/", htmlHandler)
	port := NPSConstant.WebPort

	// 启动服务器
	LogUtil.Info(fmt.Sprintf("WEB管理后台端口 :%s", port))
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		LogUtil.Error(fmt.Sprintf("WEB管理后台启动失败 :%q", err))
		log.Fatal(err)
	}
	fmt.Printf("WEB管理后台端口 :%s\n", port)
}
