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
func (p *InvoiceLineService) FindByInvoiceId(id uint) []InvoiceLine {
	return p.InvoiceLineRepository.FindByInvoiceId(id)
}
