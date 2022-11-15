package transport

import (
	"fmt"
	"net"
	"ya_rpc/common"
	"ya_rpc/pack"
	"ya_rpc/service"
)

type YaServer struct {
	Address string
	Service service.Service
	YaFunc  []func(service service.Service, params []byte) []byte
}

func (server *YaServer) Run() {
	// listen
	listener, err := net.Listen("tcp", server.Address)
	common.CheckError(err, "ListenTCP")
	println("Listening to", listener.Addr().String())
	// run
	for {
		conn, err := listener.Accept()
		common.CheckError(err, "Accept")
		go server.handleConn(conn)
	}
}

/**
以下函数为协程处理函数
*/

func (server *YaServer) handleConn(conn net.Conn) {
	defer conn.Close()
	connFrom := conn.RemoteAddr().String()
	println("Connection from  : ", connFrom)
	// 循环处理数据包
	for {
		// 获取请求数据
		request, err := pack.ParsePack(conn)
		switch err {
		case nil:
			server.handleMsg(conn, request)
		default:
			fmt.Println("Connection closed: ", connFrom)
			return
		}
	}
}

func (server *YaServer) handleMsg(conn net.Conn, request []byte) {
	// 获取methodIdx
	methodIdx, err := pack.Getuint32(&request)
	// 判断method idx是否有效
	common.CheckGoError(err, "get methodIdx")
	bad := methodIdx >= uint32(len(server.YaFunc))
	common.JudgeGoError(bad, "wrong methodIdx")
	// 按照methodIdx调用函数
	pk := server.YaFunc[methodIdx](server.Service, request)

	// 回传响应包
	_, err = conn.Write(pk)
	common.CheckGoError(err, "SendResponse")
}
