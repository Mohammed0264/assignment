// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package users

import (
	"gorm.io/gorm"
)

// Injectors from wire.go:

func InitUserApi(db *gorm.DB) UserApi {
	userRepository := ProvideUserRepository(db)
	userService := ProvideUserService(userRepository)
	userApi := ProvideUserApi(userService)
	return userApi
}
