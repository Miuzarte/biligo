# biligo

bilibili-api in golang, bro just wrapped everything.

## 总之就是封装了大部分~~常用~~我爱用的东西

~~好多东西啊什么都有这让我怎么写~~

比如动态、视频(在线人数、AI总结、视频流地址)、专栏、番剧、音频、搜索等的信息获取/接口调用
[`ReqXXX`](request.go) [`FetchXXX`](fetch.go)

还有一键解析字符串中所有的b站链接
[`ParseLink`](regexp.go)

还有扫码登录,
用迭代器实现了轮询回调 [`Login`](login.go)
([demo/login.go](https://github.com/Miuzarte/biligoDemo/blob/main/login/login.go)),
[cookie 管理](client.go)也有个很简陋的实现

还可以注册轮询视频AI总结完成的回调
[`RegisterVideoConclusion`](polling.go)

还能监听[直播间信息流](liveMsgStream.go),
主播`开播`/`下播`/`吃红SC`/`被切断`一瞬间就来通知,
快到超速了!
([demo/liveMsgStream.go](https://github.com/Miuzarte/biligoDemo/blob/main/liveMsgStream/liveMsgStream.go))

另外就是用 `text/template` 实现了为[任意结构体](types.go)的灵活[格式化](template.go),
直接丢在 [`template/type`](template/type) 里就完事了,
还做了懒加载, 太好用了, 我管你慢不慢的
(自定义模板见 [demo/customTemplate](https://github.com/Miuzarte/biligoDemo/blob/main/customTemplate/customTemplate.go))

just `go get github.com/Miuzarte/biligo`
