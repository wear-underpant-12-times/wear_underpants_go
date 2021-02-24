![](/images/head.jpeg)

[Android端](https://github.com/wear-underpant-12-times/wear_underpant_android)

## 使用方法


```
-m 模式：服务端/客户端（server/client）
-p 本地端口
-addr 服务器地址端口

Eg：
$ go run main.go -m server -p 8082  // 在8082端口启动服务器
$ go run main.go -m client -p 8082 -addr 123.123.123.123:8082     // 在8082端口启动cosk5客户端，连接远程服务器123.123.123.123:8082
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
