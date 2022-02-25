# quicServer

简单的quic-go服务器以及client实现，包含证书生成以及本地静态文件，用于测试quic的各类需求。

## Quick start

### Requirments
Go version>=1.16

### Installation
```go

git clone https://github.com/Zenger-sun/quicServer.git

// 根据需要修改配置 
vim cert/ca.conf
vim cert/server.conf

sh cert/gen.sh // 配置了默认值，执行时只需要一路回车

go mod download // 下载依赖包
go mod vendor // 将依赖复制到vendor下
```

### Build & Run

```go

go build -o server.exe
```

双击启动server.exe，如出现秒退，请检查证书是否生成  

附上benchmark

goos: windows  
goarch: amd64  
pkg: quicServer/client  
cpu: Intel(R) Core(TM) i7-7700 CPU @ 3.60GHz  
Benchmark_httpsServer  
Benchmark_httpsServer-8   	    2899	    406310 ns/op  
PASS  