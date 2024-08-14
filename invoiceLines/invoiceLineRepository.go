package invoiceLines

import (
	"fmt"
	"gorm.io/gorm"
)

type InvoiceLineRepository struct {
	Db *gorm.DB
}

func ProvideInvoiceLineRepository(db *gorm.DB) InvoiceLineRepository {
	return InvoiceLineRepository{Db: db}
}
func (p *InvoiceLineRepository) Create(invoiceLine InvoiceLine) error {
	fmt.Println("inside invoice line")
	fmt.Println(invoiceLine)
	return p.Db.Model(&InvoiceLine{}).Create(&invoiceLine).Error
}
func (p *InvoiceLineRepository) Update(invoiceLine InvoiceLine) (error, int64) {
	var count int64
	result := p.Db.Model(&InvoiceLine{}).Where("Id=?", &invoiceLine.Id).Updates(&invoiceLine).Count(&count)
	return result.Error, count
}
func (p *InvoiceLineRepository) FindByInvoiceId(id uint) []InvoiceLine {
	var invoiceLines []InvoiceLine
	p.Db.Model(&InvoiceLine{}).Where("InvoiceId=?", id).Find(&invoiceLines)
	return invoiceLines
}
