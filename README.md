# go-buildbot

编解码这种小事，犯不着每次都跑在线工具网站溜一圈。

## 用法

```
go-buildbot -run "{\"repo_url\":\"test\",\"command\":\"go build ./...\"}"
go-buildbot -list
go-buildbot -watch 60 -f config.json
```
