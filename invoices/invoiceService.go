package invoices

type InvoiceService struct {
	InvoiceRepository InvoiceRepository
}

func ProvideInvoiceService(p InvoiceRepository) InvoiceService {
	return InvoiceService{InvoiceRepository: p}
}
func (p *InvoiceService) Create(invoice Invoice) error {
	return p.InvoiceRepository.Create(invoice)

}
func (p *InvoiceService) Update(invoice Invoice) (error, int64) {
	return p.InvoiceRepository.Update(invoice)
}
func (p *InvoiceService) FindAll() []Invoice {
	return p.InvoiceRepository.FindAll()
}
func (p *InvoiceService) FindById(invoiceUniqueId string) ([]Invoice, error) {
	return p.InvoiceRepository.FindById(invoiceUniqueId)
}
func (p *InvoiceService) Delete(invoiceUniqueId string) (error, int64) {
	return p.InvoiceRepository.Delete(invoiceUniqueId)
}
