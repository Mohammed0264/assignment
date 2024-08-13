package products

type ProductDto struct {
	Id             uint    `json:"id"`
	Name           string  `json:"name" validate:"required"`
	Barcode        string  `json:"barcode" validate:"required"`
	QuantityOnHand float64 `json:"quantity_on_hand" validate:"gt=-1"`
	Price          float64 `json:"price" validate:"required,gt=0"`
	Supplier       uint    `json:"supplier" validate:"required"`
	ProductImage   string  `json:"product_image" validate:"required"`
}

// for QuantityOnHand i did not put required because it does not accept 0 while it is greater than -1
// for product image in production you can use validate:url instead of required this is more precise
