package contraller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"library/dao/mysql"
	"library/model"
)

// RegisterHandler 注册功能
func RegisterHandler(c *gin.Context) {
	p := new(model.User)
	//参数校验和绑定
	err := c.BindJSON(p)
	if err != nil {
		c.JSON(400, gin.H{"Message": err.Error()})
		return
	}
	u := model.User{Name: p.Name}
	row := mysql.DB.Where(&u).First(&u).Row()
	if row == nil {
		//将数据存储到数据库中
		mysql.DB.Create(p)
		c.JSON(200, gin.H{"Message": "success"})
	} else {
		c.JSON(403, gin.H{"Message": "该用户名已存在"})
	}
}

func LoginHandler(c *gin.Context) {
	p := new(model.User)
	//参数校验和绑定
	err := c.BindJSON(p)
	if err != nil {
		c.JSON(400, gin.H{"Message": err.Error()})
	}
	//检验用户及密码是否正确
	u := model.User{Name: p.Name, Password: p.Password}
	row := mysql.DB.Where(&u).First(&u).Row()
	if row == nil {
		c.JSON(403, gin.H{"Message": "用户名或密码错误"})
		return
	}
	//设置token
	token := uuid.New().String()
	mysql.DB.Model(&u).Update("token", token)
	c.JSON(200, gin.H{"Message": "success , your token is :" + token})
}

// LogoutHandler 登出/注销
func LogoutHandler(c *gin.Context) {
	p := new(model.User)
	//参数校验和绑定
	err := c.BindJSON(p)
	if err != nil {
		c.JSON(400, gin.H{"Message": err.Error()})
	}
	//检验用户及密码是否正确
	u := model.User{Name: p.Name, Password: p.Password}
	row := mysql.DB.Where(&u).First(&u).Row()
	if row == nil {
		c.JSON(403, gin.H{"Message": "用户名或密码错误"})
		return
	}
	//将token设置为空
	mysql.DB.Model(&u).Update("token", "")
	c.JSON(200, gin.H{"Message": "success"})
}
