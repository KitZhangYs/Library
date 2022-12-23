package model

type BookUser struct {
	UserID int64 `gorm:"primaryKey"`
	BookID int64 `gorm:"primaryKey"`
}

func (BookUser) TableName() string {
	return "books_users"
}
