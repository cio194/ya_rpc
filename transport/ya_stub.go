package transport

import (
	"net"
	"ya_rpc/common"
	"ya_rpc/pack"
)

type YaStub struct {
	conn net.Conn
}

func NewYaStub(address string) *YaStub {
	// 创建连接
	conn, err := net.Dial("tcp", address)
	common.CheckError(err, "dial failed")
	return &YaStub{conn}
}

func (stub *YaStub) RemoteCall(pk []byte) ([]byte, error) {
	// 发送请求
	_, err := stub.conn.Write(pk)
	if err != nil {
		return nil, err
	}
	// 接收响应
	return pack.ParsePack(stub.conn)
}
