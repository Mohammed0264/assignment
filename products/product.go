package products

import "gorm.io/gorm"

type Product struct {
	Id             uint           `gorm:"column:id; primary_key; auto_increment"`
	Name           string         `gorm:"column:name;"`
	Barcode        string         `gorm:"column:barcode; unique"`
	QuantityOnHand float64        `gorm:"column:quantity_on_hand"`
	Price          float64        `gorm:"column:price"`
	Supplier       uint           `gorm:"column:supplier"`
	ProductImage   string         `gorm:"column:product_image; default:Null1"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at; index"`
}
