package invoiceLines

func ToInvoiceLine(invoiceLineDto InvoiceLineDto) InvoiceLine {
	return InvoiceLine{Id: invoiceLineDto.Id, InvoiceId: invoiceLineDto.InvoiceId, ItemId: invoiceLineDto.ItemId,
		Quantity: invoiceLineDto.Quantity, LinePrice: invoiceLineDto.LinePrice}
}
func ToInvoiceDto(invoiceLine InvoiceLine) InvoiceLineDto {
	return InvoiceLineDto{Id: invoiceLine.Id, InvoiceId: invoiceLine.InvoiceId,
		ItemId: invoiceLine.ItemId, Quantity: invoiceLine.Quantity, LinePrice: invoiceLine.LinePrice}
}
func ToInvoiceDTOs(invoiceLine []InvoiceLine) []InvoiceLineDto {
	invoiceLineDTOs := make([]InvoiceLineDto, len(invoiceLine))
	for index, value := range invoiceLine {
		invoiceLineDTOs[index] = ToInvoiceDto(value)
	}
	return invoiceLineDTOs
}
