package suppliers

type SupplierDto struct {
	Id    uint   `json:"id"`
	Name  string `json:"name" validate:"required"`
	Phone string `json:"phone" validate:"required"`
}
