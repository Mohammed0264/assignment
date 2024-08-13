package invoices

import "gorm.io/gorm"

type Invoice struct {
	Id              uint           `gorm:"column:id;primary_key;auto_increment"`
	InvoiceUniqueId string         `gorm:"column:invoice_unique_id;unique"`
	Customer        uint           `gorm:"column:customer"`
	InvoiceTotal    float64        `gorm:"column:invoice_total"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at;index"`
}
