package invoices

import "gorm.io/gorm"

type InvoiceRepository struct {
	Db *gorm.DB
}

func ProvideInvoiceRepository(db *gorm.DB) InvoiceRepository {
	return InvoiceRepository{Db: db}
}
func (p *InvoiceRepository) Create(invoice Invoice) error {
	result := p.Db.Model(&Invoice{}).Create(&invoice)
	return result.Error
}
func (p *InvoiceRepository) Update(invoice Invoice) (error, int64) {
	var count int64
	result := p.Db.Model(&Invoice{}).Where("Id=?", &invoice.Id).Updates(&invoice).Count(&count)
	return result.Error, count
}
func (p *InvoiceRepository) FindAll() []Invoice {
	var invoices []Invoice
	p.Db.Find(&invoices)
	return invoices
}
func (p *InvoiceRepository) FindById(invoiceUniqueId string) ([]Invoice, error) {
	var invoices []Invoice
	result := p.Db.Model(&Invoice{}).Where("InvoiceUniqueId=?", &invoiceUniqueId).Find(&invoices)
	return invoices, result.Error
}
func (p *InvoiceRepository) Delete(invoiceUniqueId string) (error, int64) {
	result := p.Db.Model(&Invoice{}).Delete(&Invoice{}, "InvoiceUniqueId=?", &invoiceUniqueId)
	return result.Error, result.RowsAffected
}
