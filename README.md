## 简单的聊天应用

 - 基于 golang 构建的web服务器，集成了 websocket 实现
   - 基于 viper 实现的配置解析
   - mysql 作为用户数据存储()
   - redis 作为最近消息存储，与缓冲

 - 前端页面基于 vue 生态搭建
   - 配合 goland 的 template 模板解析
   - UI 库使用的是 element-ui
   - 网络请求基于 xhr 封装
   - 图标使用的是开源的 iconfont



### 构建环境

- `go SDK 1.20+` 版本


### 项目构建

 - 进入到项目的根目录

```shell
cd gochat
```

 - 执行构建脚本

 ```shell
 sh ./build/build.sh
 ```

