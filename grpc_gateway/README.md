# 编译
编译google.api
`protoc -I . --go_out=plugins=grpc,Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor:. google/api/*.proto`

出错

> plugins=grpc,Mgoogle/protobuf/descriptor.proto=github.com\golang\protobuf\protoc-gen-go\descriptor;./: No such file or directory


最终解决是参考官方github仓库
https://github.com/grpc-ecosystem/grpc-gateway

```bash
protoc -I/usr/local/include -I. \
  -I$GOPATH/src \
  -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
  --go_out=plugins=grpc:. \
  path/to/your_service.proto
```
稍微改编下：
```bash
protoc -I. -I$GOPATH/src -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=plugins=grpc:. hello.proto
```

----
上面这块其实是生成rpc代码，下面还是需要生成reverse-proxy代码的，也可以参考github仓库

最终命令
` protoc -I. -I$GOPATH/src -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --grpc-gateway_out=logtostderr=true:. hello.proto`

> 还有按照官方的教程，是不需要再hello.proto的同级目录再创建google/api，只需要在编译的时候指定一个从哪里找就可以了。
> 作者这样写的话，编译时只需要从当前目录找就可以了，也是挺好的

编写cmd
用的是https://github.com/spf13/cobra


---
note:
教程中：
```golang
func (h helloService) SayHelloWorld(ctx context.Context, r *pb.HelloWorldRequest) (*pb.HelloWorldResponse, error) {
	return &pb.HelloWorldResponse{
		Message : "test",
	}, nil
}
```
参数中的context，用的是`"golang.org/x/net/context"`
我用的是golang自带的包`context`

---
原作者的证书用的是自签名的证书，导致用postman不知道该如何访问。 自签名证书和单向认证的ca tls还是有区别的。
如果不需要ssl，只是http，可以直接参考官方github仓库的例子：https://github.com/grpc-ecosystem/grpc-gateway
如果是基于ca 的单向认证，我们甚至都不需要干什么就可以直接https请求了。
https://blog.csdn.net/ONS_cukuyo/article/details/79172242

> 为此才有下节，我们试着用ca tls，并且单向认证，客户端直接https访问就可以了。自签名证书客户端应该也可以直接访问

> 有点乱，有待思考：（自签名只是签名机构不同，其实和ca的原理是相同的），不能访问的原因应该是系统没有导入根证书吧。系统不信任我们自己的签名，我们自己扮演ca机构，不被系统信任。

但是如果我们的微服务只能是内部访问，那么任何客户端直接可以通过https请求的方式就行不通。只是在向外部提供HTTPS接口时这样才好。