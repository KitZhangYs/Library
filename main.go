package main

import (
	"library/dao/mysql"
	"library/router"
)

func main() {
	mysql.InitMysql()        //连接数据库
	r := router.InitRouter() //初始化路由
	err := r.Run(":8080")    //在8080端口运行
	//错误处理
	if err != nil {
		return
	}
}
