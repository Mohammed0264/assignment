package suppliers

import (
	"github.com/google/wire"
	"gorm.io/gorm"
)

func initSupplierApi(db *gorm.DB) SupplierAPI {
	wire.Build(ProvideSupplierApi, ProvideSupplierService, ProvideSupplierRepository)
	return SupplierAPI{}
}
