# A laravel broadcasting server

// todo

# Build

```
go get github.com/3tnet/laravel-broadcasting-server-go
cd $GOPATH/src/github.com/3tnet/laravel-broadcasting-server-go/cmd
go build -o server main.go
```

# Quick Start

```
./server
```
// todo


# Use docker
```
docker build -t laravel-broadcasting-server .
docker run -d -p 9999:9999 laravel-broadcasting-server
```

或者直接使用阿里云镜像
```
docker run -d -p 9999:9999 registry.cn-hangzhou.aliyuncs.com/ty666/laravel-broadcasting-server-go
```


# Usage

```
  -auth_endpoint string
        Auth endpoint
  -auth_host string
        Auth host
  -cors_allowed_origin string
        Cors header allowedOrigins
  -host string
        Laravel broadcasting server host

```
example:
```
./server -host="127.0.0.1:8081"
```

# Run in the background

see [https://zhuanlan.zhihu.com/p/21839884](https://zhuanlan.zhihu.com/p/21839884)
