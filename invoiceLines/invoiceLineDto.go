package invoiceLines

type InvoiceLineDto struct {
	Id        uint    `json:"id"`
	InvoiceId uint    `json:"invoiceId"`
	ItemName  string  `json:"itemName"`
	ItemId    uint    `json:"item_id"`
	Quantity  float64 `json:"quantity"`
	LinePrice float64 `json:"line_price"`
}
