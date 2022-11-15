package main

import (
	"fmt"
	"ya_rpc/ztest/ya_rpc"
)

func main() {
	const address = "127.0.0.1:50000"
	client := ya_rpc.NewYaClient(address)

	sum, err := client.Sum(1.1, 2.2)
	fmt.Println(sum, err)

	upper, err := client.Upper("abd")
	fmt.Println(upper, err)
}
