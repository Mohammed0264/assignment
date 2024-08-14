package invoiceLines

import (
	"github.com/google/wire"
	"gorm.io/gorm"
)

func initInvoiceLineService(db *gorm.DB) InvoiceLineService {
	wire.Build(ProvideInvoiceService, ProvideInvoiceLineRepository)
	return InvoiceLineService{}
}
