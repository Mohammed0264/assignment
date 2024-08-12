package customers

type CustomerDto struct {
	Id        uint    `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Address   string  `json:"address"`
	Phone     string  `json:"phone"`
	Balance   float64 `json:"balance" `
}
