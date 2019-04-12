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
这篇文章里讲到的中间人攻击部分是错误，是因为它的请求流程就是错的。它的client给server请求时，第一步就发送随机数了，这是不对的。随机数应该是服务端把证书给客户端返回后，客户端把随机数用公钥加密后给服务端。这样就不会有中间人攻击了，因为即使他拿到了数据，也不知道数据的具体内容是什么。
https://www.wosign.com/info/https_tls_ssl_http.htm
https://www.barretlee.com/blog/2016/04/24/detail-about-ca-and-certs/?utm_source=tuicool&utm_medium=referral
https://blog.csdn.net/ustccw/article/details/76691248

任何个体/组织都可以扮演 CA 的角色，只不过难以得到客户端的信任。浏览器默认信任的CA大厂有好几个。

> 这里插一点签名验签的理解:签名验签作用有验证身份，确保数据没有被篡改。首先你用hash函数将原始数据进行哈希得到摘要，然后用私钥签名（hash+私钥），我用对应的公钥验签，先是用公钥进行解密（公钥+签名）得到hash，然后对原始数据进行哈希，对比自己解密之后得到的数据，如果相同，验签成功。说明这个签名确实是你签名的。
> 如果攻击者修改了原始数据，没有改签名。那我用你的公钥解密签名后得到的hash就和原始数据的hash不同了。验签失败。这样就确保了数据没有被篡改。
> 如果攻击者修改了原始数据，并且用自己私钥进行签名，那我用你的公钥解密签名后还是和原始数据的hash对不上（只有用攻击者的公钥解密签名才能对上），验签失败。这样既能确保数据没有被篡改，也能确保签名着的身份。
> 如果攻击者修改了原始数据，并且用自己的私钥签名。并且把我本地存储的你的公钥也悄悄替换成他的了。那我验签成功了。 攻击成功了，怎么办？
这就需要CA证书的作用了。你拿你的公钥到CA认证下得到证书。以后给我发送消息时，把证书也附加上。我用ca根证书就可以验证证书的完整性。
这样攻击者就真的没有办法了。
