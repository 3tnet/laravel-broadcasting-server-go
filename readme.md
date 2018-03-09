# A laravel broadcasting server

[ English | [中文](https://github.com/3tnet/laravel-broadcasting-server-go/blob/master/readme-zh_CN.md "中文") ]

# Build

```
go get github.com/3tnet/laravel-broadcasting-server-go
cd $GOPATH/src/github.com/3tnet/laravel-broadcasting-server-go/cmd
go build -o server main.go
```
# Download

[download v0.1](https://github.com/3tnet/laravel-broadcasting-server-go/releases/tag/v0.1 "download v0.1")

# Quick Start

```
./server
```
// todo


# Using with docker
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

# config laravel
## using Pusher client
If you are broadcasting your events over Pusher, you should install the Pusher PHP SDK using the Composer package manager:
```
composer require pusher/pusher-php-server "~3.0"
```

Then add `laravel broadcasting server's` `host` and `port` and set `app id` (don't forget it!)
 `secret` not required, into `config/broadcasting.php` file `pusher` node

```
'pusher' => [
    'driver' => 'pusher',
    'key' => env('PUSHER_APP_KEY'),
    'secret' => null,
    'app_id' => env('PUSHER_APP_ID'),
    'options' => [
        'host' => 'localhost',
        'port' => 9999,
    ],
],
```

## Using Redis
todo


## Frontend
Install socket.io-client (Currently only support 1.x version)
```
npm install socket.io-client@1.x
```

```
window.io = require('socket.io-client');
import Echo from "laravel-echo"

window.Echo = new Echo({
    broadcaster: 'socket.io',
    host: window.location.hostname + ':6001'
});
```
See https://laravel.com/docs/5.6/broadcasting