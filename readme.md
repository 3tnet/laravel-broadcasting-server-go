# A laravel broadcasting server

// todo

# Build

```
go get -u
go build -o server cmd/main.go
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
  -host string
        Image Server Host (default ":9999")
```
example:
```
./server -host="127.0.0.1:8081"
```

# Run in the background

see [https://zhuanlan.zhihu.com/p/21839884](https://zhuanlan.zhihu.com/p/21839884)
