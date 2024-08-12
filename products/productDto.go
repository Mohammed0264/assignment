package products

type ProductDto struct {
	Id             uint    `json:"id"`
	Name           string  `json:"name"`
	Barcode        string  `json:"barcode"`
	QuantityOnHand float64 `json:"quantity_on_hand"`
	Price          float64 `json:"price"`
	Supplier       uint    `json:"supplier"`
	ProductImage   string  `json:"product_image"`
}
