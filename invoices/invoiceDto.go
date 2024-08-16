package invoices

import (
	"assignment/invoiceLines"
	"time"
)

type InvoiceDto struct {
	Id              uint      `json:"id"`
	InvoiceUniqueId string    `json:"uniqueInvoiceId"`
	InvoiceDate     time.Time `json:"invoiceDate"`
	Customer        uint      `json:"customer"`
	InvoiceTotal    float64   `json:"invoiceTotal"`
}

type InvoiceReceiver struct {
	Id              uint                          `json:"id"`
	InvoiceUniqueId string                        `json:"uniqueInvoiceId"`
	Customer        uint                          `json:"customer" validate:"required"`
	InvoiceDate     time.Time                     `json:"invoiceDate"`
	InvoiceTotal    float64                       `json:"invoiceTotal"`
	InvoiceLine     []invoiceLines.InvoiceLineDto `json:"line"`
}
type InvoiceSender struct {
	Id              uint                          `json:"id"`
	InvoiceUniqueId string                        `json:"uniqueInvoiceId"`
	Customer        uint                          `json:"customer"`
	CustomerName    string                        `json:"customerName"`
	InvoiceDate     time.Time                     `json:"invoiceDate"`
	InvoiceTotal    float64                       `json:"invoiceTotal"`
	InvoiceLineDto  []invoiceLines.InvoiceLineDto `json:"invoiceLine"`
}
type InvoiceUpdate struct {
	Id                uint                          `json:"id" validate:"required"`
	InvoiceUniqueId   string                        `json:"uniqueInvoiceId" validate:"required"`
	OriginalCustomer  uint                          `json:"customer" validate:"required"`
	UpdateCustomer    uint                          `json:"update_customer" validate:"required"`
	InvoiceDate       time.Time                     `json:"invoiceDate" validate:"required"`
	InvoiceTotal      float64                       `json:"invoiceTotal"`
	UpdateInvoiceLine []invoiceLines.InvoiceLineDto `json:"invoiceLineUpdate"`
	InvoiceLineDto    []invoiceLines.InvoiceLineDto `json:"invoiceLine"`
}
