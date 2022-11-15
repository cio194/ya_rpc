package ya_rpc

import (
	"ya_rpc/common"
	"ya_rpc/pack"
	"ya_rpc/service"
	"ya_rpc/transport"
)

/**
YaClient
*/

type YaClient struct {
	stub *transport.YaStub
}

func NewYaClient(address string) *YaClient {
	return &YaClient{stub: transport.NewYaStub(address)}
}

func (client *YaClient) Sum(arg0 float64, arg1 float64, arg2 string) (float64, string, error) {
	// 计算request data长度
	requestLen := uint32(4 + 8 + 8 + uint32(4+len(arg2)))

	// 数据包打包：dataLen request_data(methodIdx args)
	pk := pack.NewPack(requestLen)
	request := pk[4:]
	pack.Putuint32(&request, 0)
	pack.Putfloat64(&request, arg0)
	pack.Putfloat64(&request, arg1)
	pack.Putstring(&request, arg2)

	// 传输（远程调用）
	response, err := client.stub.RemoteCall(pk)
	common.CheckError(err, "RemoteCall")

	// 响应解析
	ret0, err := pack.Getfloat64(&response)
	common.CheckError(err, "Wrong response")
	ret1, err := pack.Getstring(&response)
	common.CheckError(err, "Wrong response")
	errRet, err := pack.Geterror(&response)
	common.CheckError(err, "Wrong response")
	return ret0, ret1, errRet
}

/**
YaServer
*/

func NewYaServer(address string, ser service.Service) *transport.YaServer {
	var yaFunc [1]func(ser service.Service, params []byte) []byte
	yaFunc[0] = yaSum

	return &transport.YaServer{Address: address, Service: ser, YaFunc: yaFunc[:]}
}

func yaSum(ser service.Service, params []byte) []byte {
	// 获取参数
	arg0, err := pack.Getfloat64(&params)
	common.CheckGoError(err, "Wrong request")
	arg1, err := pack.Getfloat64(&params)
	common.CheckGoError(err, "Wrong request")
	arg2, err := pack.Getstring(&params)
	common.CheckGoError(err, "Wrong request")

	// 执行函数
	ret0, ret1, err := ser.Sum(arg0, arg1, arg2)

	// 包装响应：dataLen data(rets err)
	responseLen := uint32(8 + uint32(4+len(ret1)) + 4 + pack.ErrLen(err))
	pk := pack.NewPack(responseLen)
	data := pk[4:]
	pack.Putfloat64(&data, ret0)
	pack.Putstring(&data, ret1)
	pack.Puterror(&data, err)

	return pk
}
