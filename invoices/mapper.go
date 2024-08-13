package invoices

func ToInvoice(invoiceDto InvoiceDto) Invoice {
	return Invoice{Id: invoiceDto.Id, InvoiceUniqueId: invoiceDto.InvoiceUniqueId, Customer: invoiceDto.Customer,
		InvoiceTotal: invoiceDto.InvoiceTotal}
}
func ToInvoiceDto(invoice Invoice) InvoiceDto {
	return InvoiceDto{Id: invoice.Id, InvoiceUniqueId: invoice.InvoiceUniqueId, Customer: invoice.Customer,
		InvoiceTotal: invoice.InvoiceTotal}
}
func ToInvoiceDTOs(invoices []Invoice) []InvoiceDto {
	invoiceDTOs := make([]InvoiceDto, len(invoices))
	for index, value := range invoices {
		invoiceDTOs[index] = ToInvoiceDto(value)
	}
	return invoiceDTOs
}
