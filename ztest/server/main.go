package main

import (
	"strings"
	"ya_rpc/ztest/ya_rpc"
)

type ServiceImpl struct {
}

func (service *ServiceImpl) Sum(a float64, b float64) (float64, error) {
	return a + b, nil
}

func (service *ServiceImpl) Upper(s string) (string, error) {
	return strings.ToUpper(s), nil
}

func main() {
	const address = "0.0.0.0:50000"
	server := ya_rpc.NewYaServer(address, &ServiceImpl{})
	server.Run()
}
