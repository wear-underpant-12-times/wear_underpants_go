![](/images/head.jpeg)

## 使用方法

[下载](https://github.com/wear-underpant-12-times/wear_underpants_go/releases/tag/0.0.11)对应操作系统的二进制包。
### 客户端
```
$ ./client.exe -h
Usage of C:\Users\Administrator\Desktop\mycode\wear_underpants_go\client\client.exe:
  -addr string    # 服务器地址，格式 ip:port
        remote addr:port (default "127.0.0.1:8082")
  -h    this help
  -p string       # 本地socks5端口
        local socks5 port (default "8081")


```
### 服务器
```
[root@23 server]# ./server -h
Usage of ./server:
  -h    this help
  -p string # 端口号
        port (default "8082")
```

## 协议

### 1.握手

|      1byte      |     0~255byte    |
|-----------------|------------------|
|  地址长度(字节)  | aes(base64(地址)) |


### 2.通信

|   2byte   |    0~65535byte   |
|-----------|------------------|
|  数据长度  | aes(base64(数据))|
