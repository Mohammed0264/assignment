package invoiceLines

type InvoiceLineDto struct {
	Id        uint    `json:"id"`
	InvoiceId uint    `json:"invoiceId" validate:"required"`
	ItemName  string  `json:"itemName"`
	ItemId    uint    `json:"item_id" validate:"required"`
	Quantity  float64 `json:"quantity" validate:"required"`
	LinePrice float64 `json:"line_price" validate:"required"`
}
