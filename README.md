# go-buildbot

简易 CI 构建工具。执行命令并记录结果，支持持续监控。

## 用法

```
go-buildbot -run "{\"repo_url\":\"test\",\"command\":\"go build ./...\"}"
go-buildbot -list
go-buildbot -watch 60 -f config.json
```
