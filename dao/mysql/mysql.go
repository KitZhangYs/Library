package mysql

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"library/model"
)

var DB *gorm.DB

func InitMysql() {
	dsn := "root:@tcp(localhost:3306)/library?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}
	DB = db
	DB.AutoMigrate(model.Book{}, model.User{})
}
