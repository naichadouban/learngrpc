用grpc测试tls，疑惑基本接触，但是没有实验成功，我们有用golang完成实验。

在server.go中，我们起了一个https服务。
client.go
```go

func main() {
	client := &http.Client{}
	resp, err := client.Get("https://localhost:8012/hello")
	if err != nil {
		panic("failed to connect: " + err.Error())
	}
	content, _ := ioutil.ReadAll(resp.Body)
	s := strings.TrimSpace(string(content))

	fmt.Println(s)
}
```
然后运行client,报错
```go
panic: failed to connect: Get https://localhost:8012/hello: x509: certificate signed by unknown authority

goroutine 1 [running]:
main.main()
	/home/xu/gopath/src/github.com/naichadouban/learngrpc/golang-tls/client.go:14 +0x1eb
exit status 2
```

# 解决方法1
我们可以忽略验证
```go
func main() {
	//roots := x509.NewCertPool()
	//ok := roots.AppendCertsFromPEM([]byte(rootPEM))
	//if !ok {
	//	panic("failed to parse root certificate")
	//}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify:true,
		},
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Get("https://localhost:8012/hello")
	if err != nil {
		panic("failed to connect: " + err.Error())
	}
	content, _ := ioutil.ReadAll(resp.Body)
	s := strings.TrimSpace(string(content))

	fmt.Println(s)
}

```
可以成功。但是这明显不好，点开 InsecureSkipVerify的注释说明，可以会有中间人攻击。
# 解决办法2:
把证书放在client代码中
```go

func main() {
	roots := x509.NewCertPool()
	pem, err := ioutil.ReadFile("./conf/server.pem")
	if err != nil{
		log.Printf("read crt file error:%v\n",err)
	}
	ok := roots.AppendCertsFromPEM(pem)
	if !ok {
		panic("failed to parse root certificate")
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:roots,
		},
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Get("https://localhost:8012/hello")
	if err != nil {
		panic("failed to connect: " + err.Error())
	}
	content, _ := ioutil.ReadAll(resp.Body)
	s := strings.TrimSpace(string(content))

	fmt.Println(s)
}
```
这样也是成功的。
上面这两中都是预料中的

# 方法3
把证书导入到系统中。为了方便，我们在ubuntu系统上进行的测试。

ubuntu很简单，我们只需要把server.pem复制到/etc/ssl/certs/下面就可以了。
然后用这段就可以了。
```go

func main() {
        client := &http.Client{}
        resp, err := client.Get("https://localhost:8012/hello")
        if err != nil {
                panic("failed to connect: " + err.Error())
        }
        content, _ := ioutil.ReadAll(resp.Body)
        s := strings.TrimSpace(string(content))

        fmt.Println(s)
}

```
这符合猜想，只是在grpc中没能完成这个实验，为什么呢？
# 总结

那为什么，在golang中我们把证书添加到本地就可以。grpc中就不可以呢？
http://singlecool.com/2017/08/18/golang-https/

> 我们知道浏览器保存了一个常用的CA证书列表，那么用golang的HttpClient请求HTTPS的服务时，它所受信任的证书列表在哪里呢？
>查看golang的源码发现在目录src/crypto/x509下有针对各个操作系统获取根证书的实现，例如root_linux.go中记录了各个Linux发行版根证书的存放路径：

golang把系统中的证书当成了自己信任列表的一部人。所以我们把证书添加到系统的路径中，golang就信任了。

但grpc没有把系统证书路径加入自己信任列表这个过程。

在demo4-stream-tls，我们也看到有人在讨论这个问题。还没有很好的方案。可能需求小众，不是很强烈。