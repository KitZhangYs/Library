package middleware

import (
	"github.com/gin-gonic/gin"
	"library/dao/mysql"
	"library/model"
)

// AuthMiddleware 权限认定，若token错误则不可以进行对图书的操作
func AuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		//从请求头中获取token
		token := c.Request.Header.Get("token")
		//新建User模型，并从users表中查询该token是否存在
		var u model.User
		if rows := mysql.DB.Where("token = ?", token).First(&u).RowsAffected; rows != 1 {
			c.JSON(403, gin.H{"msg": "token 错误"})
			//禁止进行后续操作
			c.Abort()
			return
		} else if token == "" {
			c.JSON(403, gin.H{"msg": "token 错误"})
			//禁止进行后续操作
			c.Abort()
			return
		}
		//设置UserId(不知道干什么的能不动就不动)，继续进行下一步操作
		c.Set("UserId", u.Id)
		c.Next()
	}
}
