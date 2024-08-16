package customers

type CustomerDto struct {
	Id        uint    `json:"id"`
	FirstName string  `json:"first_name" validate:"required"`
	LastName  string  `json:"last_name" validate:"required"`
	Address   string  `json:"address" validate:"required"`
	Phone     string  `json:"phone" validate:"required"`
	Balance   float64 `json:"balance" validate:"required,gt=-1"`
}
