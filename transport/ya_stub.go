package transport

import (
	"errors"
	"fmt"
	"net"
	"os"
	"time"
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
	// 重试次数
	const try = 3
	// timeout (second)
	const deadline = 5

	var err error
	var data []byte
	for i := 0; i <= try; i++ {
		if i != 0 {
			fmt.Println("remote_call retry", i)
		}
		// 设置时限
		if err = SetDeadline(stub.conn, deadline); err != nil {
			return nil, err
		}
		// 发送请求
		_, err = stub.conn.Write(pk)
		if err != nil {
			if errors.Is(err, os.ErrDeadlineExceeded) {
				continue
			} else {
				return nil, err
			}
		}
		// 接收响应
		data, err = pack.ParsePack(stub.conn)
		if err != nil {
			if errors.Is(err, os.ErrDeadlineExceeded) {
				continue
			} else {
				return nil, err
			}
		} else {
			return data, nil
		}
	}
	// 超时或其他错误
	if err != nil {
		return nil, err
	}
	return data, nil
}

func SetDeadline(conn net.Conn, seconds int) error {
	timeOut := time.Now().Add(time.Duration(seconds * 1e9))
	return conn.SetDeadline(timeOut)
}
