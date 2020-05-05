## 客户端
```
go run main.go -m client -p 8082 -addr 127.0.0.1:8083
```

## 服务端
```
go run main.go -m server -p 8083
```

## android端打包
```
gomobile bind ./mobile
```