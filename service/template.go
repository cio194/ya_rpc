package service

import (
	"bytes"
	"reflect"
	"strconv"
	"ya_rpc/common"
)

// ArgList arg0 type0, arg1 type1, ...
func ArgList(ts []reflect.Type) string {
	if len(ts) == 0 {
		return ""
	}
	var buf bytes.Buffer
	for i := 0; i < len(ts); i++ {
		buf.WriteString("arg" + strconv.Itoa(i) + " " + ts[i].Name() + ", ")
	}
	s := buf.String()
	return s[:len(s)-2]
}

// TypeList type0, type1, ...
func TypeList(ts []reflect.Type) string {
	if len(ts) == 0 {
		return ""
	}
	var buf bytes.Buffer
	for i := 0; i < len(ts); i++ {
		buf.WriteString(ts[i].Name() + ", ")
	}
	s := buf.String()
	return s[:len(s)-2]
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

func RequestDataLen(ts []reflect.Type) string {
	// request data：methodIdx args
	if len(ts) == 0 {
		return "4"
	}
	var buf bytes.Buffer
	buf.WriteString("4 + ")
	for i := 0; i < len(ts); i++ {
		buf.WriteString(mLen(ts[i], "arg"+strconv.Itoa(i)) + " + ")
	}
	s := buf.String()
	return s[:len(s)-2]
}

func mLen(t reflect.Type, argName string) string {
	k := t.Kind()
	switch k {
	case reflect.Float64:
		return "8"
	case reflect.String:
		return "uint32(4 + len(" + argName + "))"
	default:
		common.PrintExit("error type")
	}
	return ""
}
