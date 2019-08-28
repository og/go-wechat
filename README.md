# wechat


[![GoDoc](https://godoc.org/github.com/og/gowechat?status.svg)](https://godoc.org/github.com/og/gowechat)


# 功能

- 安全的缓存微信接口数据
- 安全的提供持久化存储hook
- 获取access_token
- 长链接转短链接接口


## 接口测试号申请

https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1421137522



## 必读 access_token

微信部分接口需要 access_token ，多服务器各自获取 access_token 会导致其他服务的 access_token 失效。

所以必须调用方在中控服务器获取 access_token ，调用公众号的服务必须通过中控服务器获取 access_token。


> 例： 服务器Center负责与微信交互获取 access_token 并存入数据库保存1小时55分钟，1小时55分钟后重新获取新的 access_token
> 服务器B 服务器C 每次使用 access_token 均通过请求服务器Center 获取。


网页授权 access_token 与普通 access_token不同， 网页授权 access_token 不需要实现中控服务器。

可每次通过 https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1421140842 用户授权 code 自行获取网页授权access_token。


## 单机部署中控服务器 access_token

