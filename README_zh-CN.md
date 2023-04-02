<h2 align="center">license-tool</h2>

## 项目介绍
1. 这个项目的初衷是为`交付`类的程序提供一个`license授权许可`工具
2. 项目分为服务端及客户端，客户端会生成一串`特征码`，拿到服务端`请求生成license许可文件并将内容混淆`，然后客户端再去`反混淆license许可文件并校验文件内部特征码`

## 开发进度
 - [x] 服务端
 - [x] 客户端
 - [x] 客户端`特征码`生成
 - [x] 服务端`license文件`生成
 - [x] 服务端生成的`license文件`混淆
 - [x] 客户端反混淆`license文件`
 - [x] 客户端校验`license文件`内部的特征码
 - [ ] 客户端校验`license文件`时间
 - [ ] 客户端校验`license文件`截至当前的剩余天数
 - [ ] 日志记录
 - [ ] 服务端校验`license文件`是否有效
 - [ ] 服务端查看`license许可`list
 - [ ] 服务端`license许可信息`存储至数据库
 - [ ] 客户端封装为so和dll
 - [ ] Request API doc
 - [ ] 客户端Function API doc

## 安装启动
```bash
# clone the project
https://github.com/keington/license-tool.git

# enter the project directory
cd license-tool/server

# install dependency
go mod tidy

# develop
go run .
```

## 鸣谢
感谢[JetBrain](https://www.jetbrains.com/)提供的JetBrain全家桶授权License。

## 联系我
如有问题请在issues提问，会定期解答，或者也可以发送[邮件](mailto:keington@outlook.com)

## 许可和版权
[MIT License](https://github.com/keington/license-tool/blob/cc897613c01f6ff7d2745ae1eb7303ff15a59d1c/LICENSE_zh-CN)

Copyright (c) 2023 许怀安
