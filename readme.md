
```bash
PS D:\poem-gin> go run server.go
# command-line-arguments
.\server.go:16:7: undefined: setupRouter
```
该出错原因属于go的多文件加载问题，采用go run命令执行的时候，需要把待加载的.go文件都包含到参数里面
```bash
go run server.go router.go
```

对于数据库的链接，就算可以重载 config.toml 也需要重新初始化 database.DB 链接