
```golang
for i := 0; i < 3; i++ {
		// 而在 Server 端，由于 Client 已经设置了截止时间。Server 势必要去检测它
		// 否则如果 Client 已经结束掉了，Server 还傻傻的在那执行，这对资源是一种极大的浪费
		if ctx.Err() == context.Canceled {
			return nil, status.Errorf(codes.Canceled, "SearchService.Search canceled")
		}

		time.Sleep(1 * time.Second)
		log.Printf("time:%d\n", i)
	}
```

如果我们把服务端的循环等待变成`for i := 0; i < 5; i++`,就会在客户端看到超时了