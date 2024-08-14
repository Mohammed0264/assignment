package invoices

import (
	"assignment/invoiceLines"
	"gorm.io/gorm"
	"time"
)

type Invoice struct {
	Id              uint                       `gorm:"column:id;primary_key;auto_increment"`
	InvoiceUniqueId string                     `gorm:"column:invoice_unique_id;unique"`
	Customer        uint                       `gorm:"column:customer"`
	InvoiceTotal    float64                    `gorm:"column:invoice_total"`
	InvoiceDate     time.Time                  `gorm:"column:invoice_date; type:date;"`
	InvoiceLine     []invoiceLines.InvoiceLine `gorm:"foreignKey:InvoiceId"`
	DeletedAt       gorm.DeletedAt             `gorm:"column:deleted_at;index"`
}
