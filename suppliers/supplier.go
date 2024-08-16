package suppliers

import (
	"assignment/products"
	"gorm.io/gorm"
)

type Supplier struct {
	Id        uint               `gorm:"column:id; primary_key; auto_increment;"`
	Name      string             `gorm:"column:name; not null; default:Null  size:40"`
	Phone     string             `gorm:"column:phone; not null;  size:30; unique"`
	DeletedAt gorm.DeletedAt     `gorm:"column:deleted_at; index"`
	Product   []products.Product `gorm:"foreignKey:Supplier; references:Id; "`
}
