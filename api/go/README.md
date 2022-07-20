# go

## 接口

```sh
# 更新token
curl 127.0.0.1:8000/token?t=123456
# 上传短信验证码
curl 127.0.0.1:8000/code?c=1234456
```

## 运行

```sh

-c 指定配置文件
-r 是否开始自动接口检测，默认不开启

go build -o flow-check main.go
./flow-check -c config.json -r true
exe
```