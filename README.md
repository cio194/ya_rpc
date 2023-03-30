### 介绍
参考gRPC的使用方式实现的简易RPC。
* pack：请求、响应包解析模块；
* service：用户与service.go中定义远程调用接口，然后使用generator生成远程调用客户端、服务端代码；
* transport：配合generator形成远程调用客户端、服务端，它将与用户定义的远程调用接口无关的客户端、服务端代码抽离出来，以简化代码生成器的模板文件。


### 使用
以下命令以项目根目录为工作目录
1. 在server/service.go文件内，于Service接口中自定义接口函数
2. ```go run ./service/generator/main.go```，生成代码（默认生成ztest/ya_rpc/service.ya_rpc.go）
3. 在ztest/server/main.go中编写Service接口实现
4. ```go run ./ztest/server/main.go```
5. 在ztest/client/main.go中编写测试函数
6. ```go run ./ztest/client/main.go```
