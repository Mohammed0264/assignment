package invoiceLines

type InvoiceLineDto struct {
	Id        uint    `json:"id"`
	InvoiceId string  `json:"invoiceId"`
	ItemId    uint    `json:"item_id"`
	Quantity  float64 `json:"quantity"`
	LinePrice float64 `json:"line_price"`
}
