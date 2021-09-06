package main

import (
	"fmt"
	"geektime_homework_error/global"
	"geektime_homework_error/initialize"
	"geektime_homework_error/util/docker"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	mysqlDockerCli, err := docker.NewMysqlDockerCli()

	if err != nil {
		panic(err)
	}

	err = mysqlDockerCli.Start()
	if err != nil {
		panic(err)
	}

	defer func() {
		fmt.Println("执行")
		global.DB.Close()
		err = mysqlDockerCli.Remove()
		if err != nil {
			panic(err)
		}
	}()

	initialize.InitDB()

	go func() {
		r := gin.Default()
		r.GET("/todo/:id", func(c *gin.Context) {
			id := c.Param("id")
			var todo string
			err := global.DB.QueryRow("SELECT content FROM todos WHERE id=?", id).Scan(&todo)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"data":    nil,
					"code":    1,
					"message": fmt.Sprintf("找不到id为%s的todo", id),
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"code": 0,
				"data": todo,
			})
		})
		r.Run()
	}()

	time.Sleep(10 * time.Second) // 如果建表实际失败了，请适当延长时间
	global.DB.Exec("CREATE TABLE `todos` (`id` int(10) unsigned NOT NULL AUTO_INCREMENT, `content` varchar(255) NOT NULL, PRIMARY KEY (`id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;")
	global.DB.Exec("INSERT INTO todos (content) VALUES ('完成毛老师的作业');")
	fmt.Println("Todos 表新建完成")

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
