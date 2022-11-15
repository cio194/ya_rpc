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


func (client *YaClient) Sum(arg0 float64, arg1 float64, ) (float64, error) {
    // 计算request data长度
    requestLen := uint32(4 + 8 + 8)

    // 数据包打包：dataLen request_data(methodIdx args)
    pk := pack.NewPack(requestLen)
    request := pk[4:]
    pack.Putuint32(&request, 0)
    pack.Putfloat64(&request, arg0)
    pack.Putfloat64(&request, arg1)


    // 传输（远程调用）
    response, err := client.stub.RemoteCall(pk)
    common.CheckError(err, "RemoteCall")

    // 响应解析
    ret0, err := pack.Getfloat64(&response)
    common.CheckError(err, "Wrong response")

    errRet, err := pack.Geterror(&response)
    common.CheckError(err, "Wrong response")
    err = errRet

    return ret0, err
}

func (client *YaClient) Upper(arg0 string, ) (string, error) {
    // 计算request data长度
    requestLen := uint32(4 + uint32(4 + len(arg0)))

    // 数据包打包：dataLen request_data(methodIdx args)
    pk := pack.NewPack(requestLen)
    request := pk[4:]
    pack.Putuint32(&request, 1)
    pack.Putstring(&request, arg0)


    // 传输（远程调用）
    response, err := client.stub.RemoteCall(pk)
    common.CheckError(err, "RemoteCall")

    // 响应解析
    ret0, err := pack.Getstring(&response)
    common.CheckError(err, "Wrong response")

    errRet, err := pack.Geterror(&response)
    common.CheckError(err, "Wrong response")
    err = errRet

    return ret0, err
}


/**
YaServer
*/

func NewYaServer(address string, ser service.Service) *transport.YaServer {
    var yaFunc [2]func(ser service.Service, params []byte) []byte
    yaFunc[0] = yaSum
    yaFunc[1] = yaUpper

    return &transport.YaServer{Address: address, Service: ser, YaFunc: yaFunc[:]}
}


func yaSum(ser service.Service, request []byte) []byte {
    // 获取参数
    arg0, err := pack.Getfloat64(&request)
    common.CheckGoError(err, "Wrong request")
    arg1, err := pack.Getfloat64(&request)
    common.CheckGoError(err, "Wrong request")


    // 执行函数
    ret0, err := ser.Sum(arg0, arg1)

    // 包装响应：dataLen response_data(rets err)
    responseLen := uint32(8 + 4 + pack.ErrLen(err))
    pk := pack.NewPack(responseLen)
    response := pk[4:]
    pack.Putfloat64(&response, ret0)

    pack.Puterror(&response, err)

    return pk
}

func yaUpper(ser service.Service, request []byte) []byte {
    // 获取参数
    arg0, err := pack.Getstring(&request)
    common.CheckGoError(err, "Wrong request")


    // 执行函数
    ret0, err := ser.Upper(arg0)

    // 包装响应：dataLen response_data(rets err)
    responseLen := uint32(uint32(4 + len(ret0)) + 4 + pack.ErrLen(err))
    pk := pack.NewPack(responseLen)
    response := pk[4:]
    pack.Putstring(&response, ret0)

    pack.Puterror(&response, err)

    return pk
}
