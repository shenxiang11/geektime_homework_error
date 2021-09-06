package initialize

import (
	"database/sql"
	"geektime_homework_error/global"
	_ "github.com/go-sql-driver/mysql"
)

func InitDB()  {
	dsn := "docker:123456@tcp(0.0.0.0:33061)/homework?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	global.DB, err = sql.Open("mysql", dsn)

	if err != nil {
		//time.Sleep(5 * time.Second)
		//InitDB()
		panic(err)
	}
}