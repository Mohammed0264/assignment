package invoices

import (
	"github.com/google/wire"
	"gorm.io/gorm"
)

func initInvoiceApi(db *gorm.DB) InvoiceApi {
	wire.Build(ProvideInvoiceApi, ProvideInvoiceService, ProvideInvoiceRepository)
	return InvoiceApi{}
}
