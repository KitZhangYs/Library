package contraller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"library/dao/mysql"
	"library/model"
	"strconv"
)

// CreateBook 新增书籍
func CreateBook(c *gin.Context) {
	//id := c.MustGet("UserId").(string)
	//新建Book模型并从Body中获取Book信息
	p := new(model.Book)
	err := c.BindJSON(p)
	//错误处理
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	//设定状态为未借出（gorm的默认值不知为和出错了，后续再进行修正）
	p.State = "未借出"
	//入库
	mysql.DB.Create(p)
	c.JSON(200, gin.H{"Message": "新增书籍成功"})
}

// GetBookList 获取书籍清单
func GetBookList(c *gin.Context) {
	var books []model.Book
	err := mysql.DB.Find(&books)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error})
		return
	}
	c.JSON(200, gin.H{"books": books})
}

// GetBookDetail 获取书籍信息
func GetBookDetail(c *gin.Context) {
	id := c.Param("id")
	book := new(model.Book)
	err := mysql.DB.Where("id =?", id).First(book)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error})
		return
	}
	c.JSON(200, gin.H{"book": book})
}

// UpdateBook 更新书籍
func UpdateBook(c *gin.Context) {
	//新建Book模型并从Body中获取对应信息
	NewBook := new(model.Book)
	err := c.BindJSON(NewBook)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	//建立书籍Id索引
	OldBook := &model.Book{Id: NewBook.Id}
	//寻找对应书籍
	row := mysql.DB.Where(&OldBook).First(&OldBook).Row()
	if row == nil {
		c.JSON(404, gin.H{"Message": "未找到该ID"})
		return
	} else {
		//更新书籍信息
		mysql.DB.Model(&OldBook).Updates(&NewBook)
		c.JSON(200, gin.H{"NewBook": NewBook})
	}
}

// DeleteBook 删除书籍
func DeleteBook(c *gin.Context) {
	//解析书籍id
	id1 := c.Param("id")
	//格式转化
	id, _ := strconv.Atoi(id1)
	book := new(model.Book)
	//建立书籍Id索引并查找
	OldBook := &model.Book{Id: id}
	row := mysql.DB.Where(&OldBook).First(&OldBook).Row()
	if row == nil {
		c.JSON(404, gin.H{"Message": "未找到该ID"})
		return
	} else {
		mysql.DB.Where("id =?", id).First(book)
		mysql.DB.Delete(&book)
		c.JSON(200, gin.H{"Message": "删除成功"})
	}
}

// BorrowBook 借书
func BorrowBook(c *gin.Context) {
	//获取token，从users表中查找
	token := c.Request.Header.Get("token")
	var u model.User
	rows := mysql.DB.Where("token = ?", token).First(&u).RowsAffected
	if rows != 1 {
		c.JSON(403, gin.H{"msg": "token 错误"})
		return
	}
	//新建Book模型，从Body中获取Book信息
	FindBook := new(model.Book)
	err := c.BindJSON(FindBook)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	//建立BookId索引
	book := model.Book{Id: FindBook.Id}
	row := mysql.DB.Where(&book).First(&book).Row()
	if row == nil {
		c.JSON(404, gin.H{"Message": "未找到该ID"})
		return
	} else {
		//新建book to user模型，建立BookId索引
		b2u := model.BookUser{BookID: int64(book.Id)}
		//查找该Id对应的书籍
		row = mysql.DB.Where(&b2u).First(&b2u).Row()
		if row == nil {
			//设定UserId
			b2u.UserID = int64(u.Id)
			//新建借书记录
			mysql.DB.Create(b2u)
			StateChange := book
			//更改书籍信息
			StateChange.State = "已借出"
			mysql.DB.Model(&book).Updates(&StateChange)
			c.JSON(200, gin.H{"msg": "借阅成功，记得及时归还"})
		} else {
			//新建book to user模型，从books_users表中获取对应信息
			ThisB2U := new(model.BookUser)
			mysql.DB.Where("book_id =?", b2u.BookID).First(ThisB2U)
			userid := int(ThisB2U.UserID)
			ThisUser := new(model.User)
			//从users表中获取借书人的信息并发送给客户端
			mysql.DB.Where("id = ?", userid).First(ThisUser)
			msg := fmt.Sprintf("图书已被%s借走", ThisUser.Name)
			c.JSON(200, gin.H{"Message": msg})
		}
	}
}

func ReturnBook(c *gin.Context) {
	//获取token，并从users表中查找
	token := c.Request.Header.Get("token")
	var u model.User
	rows := mysql.DB.Where("token = ?", token).First(&u).RowsAffected
	if rows != 1 {
		c.JSON(403, gin.H{"msg": "token 错误"})
		return
	}
	//新建Book模型，从Body中查找
	FindBook := new(model.Book)
	err := c.BindJSON(FindBook)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	//建立BookId索引
	book := model.Book{Id: FindBook.Id}
	row := mysql.DB.Where(&book).First(&book).Row()
	if row == nil {
		c.JSON(404, gin.H{"Message": "未找到该ID"})
		return
	} else {
		//建立book to user索引
		b2u := model.BookUser{BookID: int64(book.Id)}
		row = mysql.DB.Where(&b2u).First(&b2u).Row()
		if row == nil {
			c.JSON(200, gin.H{"Message": "该图书未被借走"})
		} else {
			ThisB2U := new(model.BookUser)
			mysql.DB.Where("book_id =?", b2u.BookID).First(ThisB2U)
			userid := int(ThisB2U.UserID)
			ThisUser := new(model.User)
			mysql.DB.Where("id = ?", userid).First(ThisUser)
			//检验借书人与还书人是否对应
			if ThisUser.Id == u.Id {
				DeleteB2U := new(model.BookUser)
				mysql.DB.Where("book_id =?", ThisB2U.BookID).First(DeleteB2U)
				mysql.DB.Delete(&DeleteB2U)
				StateChange := book
				StateChange.State = "未借出"
				mysql.DB.Model(&book).Updates(&StateChange)
				c.JSON(200, gin.H{"Message": "还书成功，欢迎下次借阅"})
			}
		}
	}
}
