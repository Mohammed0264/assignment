package invoiceLines

import "gorm.io/gorm"

type InvoiceLine struct {
	Id        uint           `gorm:"column:id; primary_key; auto_increment;"`
	InvoiceId uint           `gorm:"column:InvoiceId;"`
	ItemId    uint           `gorm:"column:item_id;"`
	Quantity  float64        `gorm:"column:quantity; default:0"`
	LinePrice float64        `gorm:"column:line_price; default:0"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at; index"`
}
