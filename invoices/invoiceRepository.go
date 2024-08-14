package invoices

import (
	"fmt"
	"gorm.io/gorm"
)

type InvoiceRepository struct {
	Db *gorm.DB
}

func ProvideInvoiceRepository(db *gorm.DB) InvoiceRepository {
	return InvoiceRepository{Db: db}
}
func (p *InvoiceRepository) Create(invoice Invoice) (error, Invoice) {

	var lastCreatedInvoice Invoice
	result := p.Db.Model(&Invoice{}).Create(&invoice).Last(&lastCreatedInvoice)
	if result.Error != nil {
		fmt.Println("invoice was not created")
	}
	if result.Error == nil {
		fmt.Println("invoice was created")
	}
	fmt.Println(lastCreatedInvoice)
	return result.Error, lastCreatedInvoice
}
func (p *InvoiceRepository) Update(invoice Invoice) (error, int64) {
	var count int64
	result := p.Db.Model(&Invoice{}).Where("Id=?", &invoice.Id).Updates(&invoice).Count(&count)
	return result.Error, count
}
func (p *InvoiceRepository) FindAll() []Invoice {
	var invoices []Invoice
	result := p.Db.Model(&Invoice{}).Find(&invoices)
	if result.Error != nil {
		fmt.Println(result.Error.Error())
	}
	fmt.Println(invoices)
	return invoices
}
func (p *InvoiceRepository) FindById(invoiceUniqueId string) (Invoice, error) {
	var invoices Invoice
	result := p.Db.Model(&Invoice{}).Where("invoice_unique_id=?", &invoiceUniqueId).Find(&invoices)
	return invoices, result.Error
}
func (p *InvoiceRepository) Delete(invoiceUniqueId string) (error, int64) {
	result := p.Db.Model(&Invoice{}).Where("invoice_unique_id=?", &invoiceUniqueId).Delete(&Invoice{})
	return result.Error, result.RowsAffected
}
func (p *InvoiceRepository) FindLastInvoice() Invoice {
	var invoice Invoice
	p.Db.Model(&Invoice{}).Unscoped().Last(&invoice)
	fmt.Println("inside last")
	fmt.Println(invoice)
	return invoice
}
