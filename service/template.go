package service

import (
	"bytes"
	"reflect"
	"strconv"
	"ya_rpc/common"
)

// ArgList (arg0 type0, arg1 type1, ...,)
func ArgList(ts []reflect.Type) string {
	if len(ts) == 0 {
		return "()"
	}
	var buf bytes.Buffer
	buf.WriteByte('(')
	for i := 0; i < len(ts); i++ {
		buf.WriteString("arg" + strconv.Itoa(i) + " " + ts[i].Name() + ", ")
	}
	buf.WriteByte(')')
	return buf.String()
}

// RetList (type0, type1, ..., error)
func RetList(ts []reflect.Type) string {
	if len(ts) == 0 {
		return "error"
	}
	var buf bytes.Buffer
	buf.WriteByte('(')
	for i := 0; i < len(ts); i++ {
		buf.WriteString(ts[i].Name() + ", ")
	}
	buf.WriteString("error)")
	return buf.String()
}

// RequestLen request: methodIdx args
func RequestLen(ts []reflect.Type) string {
	if len(ts) == 0 {
		return "4"
	}
	var buf bytes.Buffer
	buf.WriteString("4 + ")
	for i := 0; i < len(ts); i++ {
		buf.WriteString(mLen(ts[i], "arg"+strconv.Itoa(i)) + " + ")
	}
	s := buf.String()
	return s[:len(s)-3]
}

// ReturnList ret0, ret1, ..., errRet
func ReturnList(ts []reflect.Type) string {
	if len(ts) == 0 {
		return "err"
	}
	var buf bytes.Buffer
	for i := 0; i < len(ts); i++ {
		buf.WriteString("ret" + strconv.Itoa(i) + ", ")
	}
	buf.WriteString("err")
	return buf.String()
}

// CallList arg0, arg1, ...
func CallList(ts []reflect.Type) string {
	if len(ts) == 0 {
		return ""
	}
	var buf bytes.Buffer
	for i := 0; i < len(ts); i++ {
		buf.WriteString("arg" + strconv.Itoa(i) + ", ")
	}
	s := buf.String()
	return s[:len(s)-2]
}

// ResponseLen response: rets err
func ResponseLen(ts []reflect.Type) string {
	if len(ts) == 0 {
		return "4 + pack.ErrLen(err)"
	}
	var buf bytes.Buffer
	for i := 0; i < len(ts); i++ {
		buf.WriteString(mLen(ts[i], "ret"+strconv.Itoa(i)) + " + ")
	}
	buf.WriteString("4 + pack.ErrLen(err)")
	return buf.String()
}

/**
下列函数类型扩展时需调整
*/

func SupportKind(k reflect.Kind) bool {
	return k == reflect.Float64 ||
		k == reflect.String ||
		k == reflect.Int64
}

func mLen(t reflect.Type, name string) string {
	k := t.Kind()
	switch k {
	case reflect.Float64:
		return "8"
	case reflect.String:
		return "uint32(4 + len(" + name + "))"
	case reflect.Int64:
		return "8"
	default:
		common.PrintExit("error type")
	}
	return ""
}

func CheckFuncLimit(serviceType reflect.Type) {
	// 无unexported方法
	// 参数支持：float64、string、int64
	// 返回值：可多返回值，类型支持与参数支持相同，尾返回值必须为error
	methodNum := 0
	for i := 0; i < serviceType.NumMethod(); i++ {
		method := serviceType.Method(i)
		if !method.IsExported() {
			common.PrintExit("no need of unexported method")
		}
		methodNum++
		// args
		for j := 0; j < method.Type.NumIn(); j++ {
			arg := method.Type.In(j)
			if !SupportKind(arg.Kind()) {
				common.PrintExit("method " + method.Name + " bad arg_type " + arg.Name())
			}
		}
		// rets
		rets := make([]reflect.Type, method.Type.NumOut())
		for j := 0; j < method.Type.NumOut(); j++ {
			rets[j] = method.Type.Out(j)
		}
		if len(rets) == 0 {
			common.PrintExit("method " + method.Name + " at least one error for return")
		}
		for j := 0; j < len(rets)-1; j++ {
			if !SupportKind(rets[j].Kind()) {
				common.PrintExit("method " + method.Name + " bad ret_type " + rets[j].Name())
			}
		}
		if rets[len(rets)-1].Name() != "error" {
			common.PrintExit("method " + method.Name + " last ret should be error")
		}
	}
	if methodNum == 0 {
		common.PrintExit("no exported method")
	}
}
