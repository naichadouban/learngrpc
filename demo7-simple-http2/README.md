现在的情况，无论是通过`client.go`访问，还是通过浏览器直接访问都是可以的。

postman还没有成功，不知道为什么。但是既然浏览器可以成功，说明现在https访问是可以的。

追加：浏览器可以成功，是我们点击了忽略验证，继续访问资源。

但是postman没有这个功能，所以postman不能访问。

```go
ClientAuth: tls.RequireAndVerifyClientCert,
ClientAuth: tls.NoClientCert,
```
这两个无论设置成什么都是可以的。


# note1：
```golang
http.ListenAndServeTLS(port, certFile, keyFile, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("accept")
		fmt.Println(r.Header.Get("Content-Type"))
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			server.ServeHTTP(w, r)
		} else {
			mux.ServeHTTP(w, r)
		}
		return
	}))
```
`server.ServeHTTP(w, r)`,这个server里面是需要https的。

我们看方法上面的注释就可以看到，要求通过HTTP2，通过tls。解决了上节的疑惑。
```go
// The provided HTTP request must have arrived on an HTTP/2
// connection. When using the Go standard library's server,
// practically this means that the Request must also have arrived
// over TLS.
```

所以我们就必须要用`http.ListenAndServeTLS`,而不是`http.ListenAndServe`

# note2
这个章节和教程中关于证书那一块也是不同的。教程中那块用的是TLS证书认证，我们用的基于CA的证书认证。


感觉还是有点乱。