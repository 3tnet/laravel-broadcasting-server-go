# A laravel broadcasting server

[ [English](https://github.com/3tnet/laravel-broadcasting-server-go/blob/master/readme.md "English") | 中文 ]

# 构建

```
go get github.com/3tnet/laravel-broadcasting-server-go
cd $GOPATH/src/github.com/3tnet/laravel-broadcasting-server-go/cmd
go build -o server main.go
```

# 快速开始

```
./server
```


# 使用 docker
```
docker build -t laravel-broadcasting-server .
docker run -d -p 9999:9999 laravel-broadcasting-server
```

或者直接使用阿里云镜像
```
docker run -d -p 9999:9999 registry.cn-hangzhou.aliyuncs.com/ty666/laravel-broadcasting-server-go
```


# 启动参数

```
  -auth_endpoint string
        Auth endpoint
  -auth_host string
        Auth host
  -cors_allowed_origin string
        Cors header allowedOrigins （设置允许跨域的网址）
  -host string
        Laravel broadcasting server host

```
例子:
```
./server -host="127.0.0.1:8081"
```

# 后台运行

see [https://zhuanlan.zhihu.com/p/21839884](https://zhuanlan.zhihu.com/p/21839884)

# 配置 laravel
## 使用 Pusher client
如果你使用 Pusher 对事件进行广播，请用 Composer 包管理器来安装 Pusher PHP SDK：
```
composer require pusher/pusher-php-server "~3.0"
```
然后在 `config/broadcasting.php` 中的 `pusher` 节点中添加 `laravel broadcasting server` 的 `host` 和 `port` 并且设置 `app id` (不要忘记！)
 `secret` 不用填。
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

## 使用 Redis
todo


## 前端
安装 socket.io-client (目前只支持 1.x 版本)
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
详见 https://laravel-china.org/docs/laravel/5.6/broadcasting