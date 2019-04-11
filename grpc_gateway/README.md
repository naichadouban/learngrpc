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