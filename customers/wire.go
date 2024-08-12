package customers

import (
	"github.com/google/wire"
	"gorm.io/gorm"
)

func initCustomerApi(db *gorm.DB) CustomerApi {
	wire.Build(ProvideCustomerApi, ProvideCustomerService, ProvideCustomerRepository)
	return CustomerApi{}
}
