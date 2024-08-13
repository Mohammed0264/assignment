package invoiceLines

type InvoiceLine struct {
	Id        uint    `gorm:"column:id; primary_key; auto_increment;"`
	InvoiceId string  `gorm:"column:invoice_id"`
	ItemId    uint    `gorm:"column:item_id;"`
	Quantity  float64 `gorm:"column:quantity; default:0"`
	LinePrice float64 `gorm:"column:line_price; default:0"`
}
