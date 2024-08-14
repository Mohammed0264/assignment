package invoices

import (
	"assignment/invoiceLines"
	"time"
)

var time1 time.Time

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
	Customer        uint                          `json:"customer"`
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
