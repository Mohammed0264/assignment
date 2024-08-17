package users

import "gorm.io/gorm"

type User struct {
	Id       uint           `gorm:"primary_key; auto_increment; column:id"`
	UserName string         `gorm:"column:user_name; not null; unique"`
	Password string         `gorm:"column:password; not null"`
	Role     string         `gorm:"column:role; not null; default:member"`
	DeleteAt gorm.DeletedAt `gorm:"column:deleted_at; index"`
}
