# 编译
编译google.api
`protoc -I . --go_out=plugins=grpc,Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor:. google/api/*.proto`

出错

> plugins=grpc,Mgoogle/protobuf/descriptor.proto=github.com\golang\protobuf\protoc-gen-go\descriptor;./: No such file or directory