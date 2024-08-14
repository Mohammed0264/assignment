// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package invoiceLines

import (
	"gorm.io/gorm"
)

// Injectors from wire.go:

func InitInvoiceLineService(db *gorm.DB) InvoiceLineService {
	invoiceLineRepository := ProvideInvoiceLineRepository(db)
	invoiceLineService := ProvideInvoiceService(invoiceLineRepository)
	return invoiceLineService
}
