# go

## 接口

```sh
# 上传短信验证码
curl 127.0.0.1:8000/code?c=1234456
# 更新token
curl 127.0.0.1:8000/token?t=123456
# 手动开启检测
curl 127.0.0.1:8000/check
# 检测是否运行
curl 127.0.0.1:8000/running
```

## 运行

```sh

-c 指定配置文件
-r 是否开始自动接口检测，默认不开启

go build -o flow-check main.go
./flow-check -c config.json -r true
```

## 日常

每天早上发一封邮件确认一下
