package main

import (
	"os"
	"reflect"
	"text/template"
	"ya_rpc/common"
	"ya_rpc/service"
)

type funcInfo struct {
	Method reflect.Method
	Args   []reflect.Type
	Rets   []reflect.Type
}

func main() {
	// 获取输出文件路径
	//flag.Parse()
	//flag.Parse()
	//if flag.NArg() != 1 {
	//	common.PrintExit("usage: outfile")
	//}
	//of := flag.Arg(0)
	of := "./ztest/ya_rpc/service.ya_rpc.go"
	outFile, err := os.OpenFile(of, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	common.CheckError(err, "open outfile "+of)

	// 检查函数限制
	// 参数类型限制、返回类型 error 限制
	serviceType := reflect.TypeOf((*service.Service)(nil)).Elem()
	service.CheckFuncLimit(serviceType)

	// 获取函数元信息
	funcInfos := GetFuncInfos(serviceType)

	// 模板注册自定义函数
	tmpl := template.New("service.tmpl")
	tmpl = RegisterFunc(tmpl)

	// 模板解析
	tmplFile := "./service/service.tmpl"
	tmpl, err = tmpl.ParseFiles(tmplFile)
	common.CheckError(err, "template ParseFiles")

	// 模板执行
	err = tmpl.Execute(outFile, funcInfos)
	common.CheckError(err, "template execute")
}

func RegisterFunc(tmpl *template.Template) *template.Template {
	funcMap := map[string]interface{}{
		"ArgList":     service.ArgList,
		"RetList":     service.RetList,
		"RequestLen":  service.RequestLen,
		"ReturnList":  service.ReturnList,
		"CallList":    service.CallList,
		"ResponseLen": service.ResponseLen,
		"ReturnError": service.ReturnError,
	}
	return tmpl.Funcs(funcMap)
}

func GetFuncInfos(serviceType reflect.Type) []funcInfo {
	funcInfos := make([]funcInfo, serviceType.NumMethod())
	for i := 0; i < serviceType.NumMethod(); i++ {
		method := serviceType.Method(i)
		funcType := method.Type
		// args
		args := make([]reflect.Type, funcType.NumIn())
		for j := 0; j < funcType.NumIn(); j++ {
			args[j] = funcType.In(j)
		}
		// rets，不包括error
		rets := make([]reflect.Type, funcType.NumOut()-1)
		for j := 0; j < funcType.NumOut()-1; j++ {
			rets[j] = funcType.Out(j)
		}
		funcInfos[i].Method = method
		funcInfos[i].Args = args
		funcInfos[i].Rets = rets
	}
	return funcInfos
}
