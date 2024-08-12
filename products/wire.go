package products

import (
	"github.com/google/wire"
	"gorm.io/gorm"
)

func initProductApi(db *gorm.DB) ProductApi {
	wire.Build(ProvideProductApi, ProvideProductService, ProvideProductRepository)
	return ProductApi{}
}
