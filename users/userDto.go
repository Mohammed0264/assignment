package users

type UserDto struct {
	Id       uint   `json:"id"`
	UserName string `json:"userName" validate:"required"`
	Password string `json:"password" validate:"required"`
}
