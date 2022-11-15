package common

import (
	"fmt"
	"os"
	"runtime"
)

func CheckError(err error, info string) {
	if err != nil {
		fmt.Println("ERROR:", info, err.Error())
		os.Exit(1)
	}
}

func CheckGoError(err error, info string) {
	if err != nil {
		fmt.Println("ERROR:", info, err.Error())
		runtime.Goexit()
	}
}

func JudgeError(bad bool, info string) {
	if !bad {
		return
	}
	// 发生错误
	fmt.Println("BAD:", info)
	os.Exit(1)
}

func JudgeGoError(bad bool, info string) {
	if !bad {
		return
	}
	// 发生错误
	fmt.Println("BAD:", info)
	os.Exit(1)
}

func PrintExit(info string) {
	fmt.Println("BAD:", info)
	os.Exit(1)
}
