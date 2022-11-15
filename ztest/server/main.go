package main

import (
	"strings"
	"ya_rpc/ztest/ya_rpc"
)

type ServiceImpl struct {
}

func (service ServiceImpl) Sum(a float64, b float64, c string) (float64, string, error) {
	return a + b, strings.ToUpper(c), nil
}

func main() {
	const address = "0.0.0.0:50000"
	server := ya_rpc.NewYaServer(address, &ServiceImpl{})
	server.Run()
}
