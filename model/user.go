package model

type User struct {
	Id       int    `gorm:"primary_key"`
	Name     string `json:"name" gorm:"not null" binding:"required"`
	Password string `json:"password" gorm:"not null" binding:"required"`
	Age      int    `json:"age" gorm:"age"`
	Token    string `json:"token"`
}
