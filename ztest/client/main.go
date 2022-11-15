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
	serviceType := reflect.TypeOf((*service.Service)(nil)).Elem()

	// todo 检查各函数是否满足限制条件
	// check

	// 通过反射获取函数元信息
	// 包装模板所需信息：函数，参数，返回值
	funcInfos := make([]funcInfo, serviceType.NumMethod())
	for i := 0; i < serviceType.NumMethod(); i++ {
		method := serviceType.Method(i)
		funcType := method.Type
		// 参数
		args := make([]reflect.Type, funcType.NumIn())
		for j := 0; j < funcType.NumIn(); j++ {
			args[j] = funcType.In(j)
		}
		// 返回值，不应包括error
		numOut := funcType.NumOut() - 1
		rets := make([]reflect.Type,numOut)
		for j := 0; j < numOut; j++ {
			rets[j] = funcType.Out(j)
		}
		funcInfos[i].Method = method
		funcInfos[i].Args = args
		funcInfos[i].Rets = rets
	}

	// 注册自定义函数
	funcMap := map[string]interface{}{
		"ArgList":        service.ArgList,
		"TypeList":       service.TypeList,
		"CallList":       service.CallList,
		"RequestDataLen": service.RequestDataLen,
	}
	tmpl := template.New("service.tmpl").Funcs(funcMap)

	// 模板生成
	var err error
	tmplFile := "/home/simple/code/ya_rpc/test/service.tmpl"
	tmpl, err = tmpl.ParseFiles(tmplFile)
	common.CheckError(err, "template ParseFiles")
	err = tmpl.Execute(os.Stdout, funcInfos)
	common.CheckError(err, "template execute")
}
