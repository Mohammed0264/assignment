package invoices

import "time"

func ToInvoice(invoiceReceiver InvoiceReceiver) Invoice {
	time1 := time.Now().Truncate(24 * time.Hour)
	return Invoice{Id: invoiceReceiver.Id, InvoiceUniqueId: invoiceReceiver.InvoiceUniqueId, Customer: invoiceReceiver.Customer,
		InvoiceTotal: invoiceReceiver.InvoiceTotal, InvoiceDate: time1}
}
func ToInvoiceDto(invoice Invoice) InvoiceDto {
	return InvoiceDto{Id: invoice.Id, InvoiceUniqueId: invoice.InvoiceUniqueId, Customer: invoice.Customer,
		InvoiceTotal: invoice.InvoiceTotal, InvoiceDate: invoice.InvoiceDate}
}
func ToInvoiceDTOs(invoices []Invoice) []InvoiceDto {
	invoiceDTOs := make([]InvoiceDto, len(invoices))
	for index, value := range invoices {
		invoiceDTOs[index] = ToInvoiceDto(value)
	}
	return invoiceDTOs
}
