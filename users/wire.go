package users

import (
	"github.com/google/wire"
	"gorm.io/gorm"
)

func initUserApi(db *gorm.DB) UserApi {
	wire.Build(ProvideUserApi, ProvideUserService, ProvideUserRepository)
	return UserApi{}
}
