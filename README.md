> 我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？

从 database/sql 源码中找到的注释
> ErrNoRows is returned by Scan when QueryRow doesn't return a row. In such a case, QueryRow returns a placeholder *Row value that defers this error until a Scan.

大致意思是：
QueryRow 执行一个查询，预计最多返回一条记录。
QueryRow 总是返回一个非零的值。错误被推迟到 行的扫描方法被调用。
如果查询没有选择任何行，*Row 的 Scan 将返回 ErrNoRows。 
否则，*Row's Scan 会扫描第一条被选中的记录，并丢弃 其余的。


我写了一个简单的 gin 应用，
先启动本地 docker app，
再运行 main.go,
顺利的话，会在 8080 端口启动，数据库也会跑在 docker 中。

数据表 todos 大概是下面这样：

|  id   | content  |
|  ----  | ----  |
| 1  | 完成毛老师的作业 |

正常情况的响应：
```
Request URL: http://localhost:8080/todo/1
Request Method: GET
Status Code: 200 OK
Response: {"code":0,"data":"完成毛老师的作业"}
```

遇到 sql.ErrNoRows 时的响应：
```
Request URL: http://localhost:8080/todo/100
Request Method: GET
Status Code: 404 Not Found
Response: {"code":1,"data":null,"message":"找不到id为100的todo"}
```

遇到 sql.ErrNoRows 时，我认为需要把这个错误 Wrap，并抛给上层。作为底层服务，我认为要告诉调用方"没有找到相应的资源"。

在我的例子中，遇到 sql.ErrNoRows 时， 我返回 404 也相当于 Wrap 了这个错误，并抛给了前端，由前端去根据 http status 或者 Response 中的 code，去显示 404 页，或者其他处理。
如果返回前端 200 的状态，但是 data 却是 null (即 todo)，前端会有二义性，这条 todo 是"空" 的 或者 没有这条 todo，此时需要双方提前约定，才能消除这种二义性，如果人员变动大的工作环境，返回200，把错误"吞掉"，我认为不合适。

而且如果做了监控，我们感知这种 404 错误，可以作出及时的处理。

---

这次遇到的其他问题

我用 go 启动 docker 的mysql，然后做连接，不 sleep 一会，会建不了表？
老师有什么优雅的方式解决这个问题吗？

前面一开始我用 gorm 试了一下，连接 db 前都得 sleep 一会

err 是
```
panic: dial tcp 127.0.0.1:33061: connect: connection refused
```

```go
initialize.InitDB()
xxxx
time.Sleep(10 * time.Second) // 如果建表实际失败了，请适当延长时间
global.DB.Exec("CREATE TABLE `todos` (`id` int(10) unsigned NOT NULL AUTO_INCREMENT, `content` varchar(255) NOT NULL, PRIMARY KEY (`id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;")
global.DB.Exec("INSERT INTO todos (content) VALUES ('完成毛老师的作业');")
fmt.Println("Todos 表新建完成")
```