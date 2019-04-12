# learngrpc
grpc学习记录


1. 带入gRPC：gRPC及相关介绍
2. 带入gRPC：gRPC Client and Server
3. 带入gRPC：gRPC Streaming, Client and Server
4. 带入gRPC：TLS 证书认证
5. 带入gRPC：基于 CA 的 TLS 证书认证
6. 带入gRPC：Unary and Stream interceptor
7. 带入gRPC：让你的服务同时提供 HTTP 接口
8. 带入gRPC：对 RPC 方法做自定义认证
9. 带入gRPC：gRPC Deadlines
10. 带入gRPC：分布式链路追踪 gRPC+Opentracing+Zipkin

# 参考文章
https://github.com/EDDYCJY/blog

其中有些地方和参考文章不同，但最终实现效果是一样的。

# 比较好的参考
grpc入门必读：gRPC Go: Beyond the basics  

https://blog.gopheracademy.com/advent-2017/go-grpc-beyond-basics/


这个人的几篇文章都值得参考下
[Tutorial, Part 1] How to develop Go gRPC microservice with HTTP/REST endpoint, middleware, Kubernetes deployment, etc.
https://medium.com/@amsokol.com/tutorial-how-to-develop-go-grpc-microservice-with-http-rest-endpoint-middleware-kubernetes-daebb36a97e9

# 回顾
方便自己回忆，不构成参考
## grpc影响
服务端主要代码
```golang
lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
if err != nil {
        log.Fatalf("failed to listen: %v", err)
}

grpcServer := grpc.NewServer()
...
pb.RegisterSearchServer(grpcServer, &SearchServer{})
grpcServer.Serve(lis)
```
客户端主要代码
```golang
var opts []grpc.DialOption
...
conn, err := grpc.Dial(*serverAddr, opts...)
if err != nil {
    log.Fatalf("fail to dial: %v", err)
}

defer conn.Close()
client := pb.NewSearchClient(conn)
```

## Simple RPC
代码在demo2
1. protoc 的简单命令要掌握，protoc help看下就差不多了

## 流式 grpc
代码在demo3-stream

流式（stream）适用场景：
1. 实时
2. 大规模数据

具体是选择流式grpc还是普通grpc要根据自己的实际业务来决定

## 给rpc加上tls证书认证
代码主要集中在demo4-stream-tls

> 在看这部分代码之前，强烈建议先读下面文章
https://www.barretlee.com/blog/2015/10/05/how-to-build-a-https-server/
https://www.barretlee.com/blog/2016/04/24/detail-about-ca-and-certs/?utm_source=tuicool&utm_medium=referral
https://blog.csdn.net/ustccw/article/details/76691248

任何个体/组织都可以扮演 CA 的角色，只不过难以得到客户端的信任。浏览器默认信任的CA大厂有好几个。

HTTPS请求
1. client->server random-client(客户端随机数)+say hello
2. server ->