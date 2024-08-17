package invoiceLines

type InvoiceLineService struct {
	InvoiceLineRepository InvoiceLineRepository
}

func ProvideInvoiceService(p InvoiceLineRepository) InvoiceLineService {
	return InvoiceLineService{InvoiceLineRepository: p}
}
func (p *InvoiceLineService) Create(invoiceLineDto InvoiceLineDto) error {

	err := p.InvoiceLineRepository.Create(ToInvoiceLine(invoiceLineDto))
	if err != nil {
		return err
	}
	return nil
}
func (p *InvoiceLineService) Update(invoiceLineDto InvoiceLineDto) (error, int64) {

	result, count := p.InvoiceLineRepository.Update(ToInvoiceLine(invoiceLineDto))
	if result != nil {
		return result, 0
	}
	return nil, count
}
func (p *InvoiceLineService) FindByInvoiceId(id uint) []InvoiceLine {
	return p.InvoiceLineRepository.FindByInvoiceId(id)
}
func (p *InvoiceLineService) Delete(id uint) (error, int64) {
	err, count := p.InvoiceLineRepository.Delete(id)
	return err, count
}
