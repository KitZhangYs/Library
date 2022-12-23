package model

type Book struct {
	Id    int    `gorm:"primary_key"`
	Name  string `json:"name" gorm:"not null" binding:"required"`
	State string `gorm:"state" gorm:"default:'Not Lented'"`
	User  []User `gorm:"many2many:books_users"`
}
