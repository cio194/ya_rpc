package main

import (
	"fmt"
	"ya_rpc/test/service"
)

func main() {
	c := service.NewSummer()
	fmt.Println(c.Sum(1.1, 2.2))
}
