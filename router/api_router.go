package router

import (
	"github.com/gin-gonic/gin"
	"library/contraller"
	"library/middleware"
)

func LoadApiRouter(r *gin.Engine) {
	r.POST("/register", contraller.RegisterHandler) //用户注册
	r.POST("/login", contraller.LoginHandler)       //用户登录
	//创建分组路由lib
	v1 := r.Group("/lib/v1")
	v1.Use(middleware.AuthMiddleware())            //权限检验
	v1.POST("/book", contraller.CreateBook)        //创建书籍
	v1.GET("/book", contraller.GetBookList)        //获得书籍清单
	v1.GET("/book/:id", contraller.GetBookDetail)  //获得具体书籍
	v1.PUT("/book", contraller.UpdateBook)         //更新书籍信息
	v1.DELETE("/book/:id", contraller.DeleteBook)  //删除书籍
	v1.POST("/book/borrow", contraller.BorrowBook) //借书功能
	v1.POST("/book/return", contraller.ReturnBook) //还书功能
}
