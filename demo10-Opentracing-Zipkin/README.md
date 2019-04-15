note：

起初没有配置ca证书，只是简单的弄了下tls认证，然后zipkin那边就一直收不到信息。
后来改成了ca tls认证，zipkin才收到了信息。


但是个人觉得应该是没有关系的啊

为此多加了一个项目demo10-Opentracing-Zipkin2

在这个项目中，我们不加认证了，试试情况