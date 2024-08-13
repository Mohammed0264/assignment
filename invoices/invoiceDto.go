package invoices

type InvoiceDto struct {
	Id              uint    `json:"id"`
	InvoiceUniqueId string  `json:"uniqueInvoiceId"`
	Customer        uint    `json:"customer"`
	InvoiceTotal    float64 `json:"invoiceTotal"`
}
