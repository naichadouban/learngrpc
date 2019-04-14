tls证书生成：
```bash
openssl ecparam -genkey -name secp384r1 -out server.key

openssl req -new -x509 -sha256 -key server.key -out server.pem -days 3650

```
后来发现ubuntu导入证书时，要`crt`格式。对比一下两种格式的生成。

可参考：https://www.barretlee.com/blog/2015/10/05/how-to-build-a-https-server/

# 问题1
为什么客户端也要传证书呢？
https://forfuncsake.github.io/post/2017/08/trust-extra-ca-cert-in-go-app/
这篇文章，指出了自签证书会存在的问题（我们这个就是自签），
已经对应的解决办法。我们可以一一试试。

首先要考虑的问题，HTTP2必须要tls吗？
https://daniel.haxx.se/blog/2015/03/06/tls-in-http2/
这篇文章有解释，大概意思就是草案并没有要求，
但是各个浏览器厂家只实现了基于tls的http2，
那如果我们的应用时基于http2的，但是又不是tls，
在他们那里就通不过了。

那grpc必须要tls吗？
不是，后面的章节有验证。

# 尝试1
直接在客户端取消验证
```go
	//c, err := credentials.NewClientTLSFromFile("../conf/server.pem", "localhost")
	//if err != nil {
	//	log.Panicln(err)
	//}
	//cc, err := grpc.Dial(port, grpc.WithTransportCredentials(c))
	cc, err := grpc.Dial(port, grpc.WithInsecure())
```
报错

> D:\gopath\src\github.com\naichadouban\learngrpc\demo4-stream-tls\client>go run stream_client.go
2019/04/14 21:27:25 rpc error: code = Unavailable desc = all SubConns are in TransientFailure, latest connection error: connection error: desc = "transport: Error while dialing dial tcp: address https://127.0.0.1:8011: too many colons in address"
panic: runtime error: invalid memory address or nil pointer dereference
[signal 0xc0000005 code=0x0 addr=0x30 pc=0x78ec72]

github这个issue
https://github.com/grpc/grpc-go/issues/2443

有人觉得可能有两种原因：
1. 服务端希望验证，客户端不验证 WithInsecure
2. 客户端证书配置错误

# 尝试2
把自签的证书添加到系统的证书中
