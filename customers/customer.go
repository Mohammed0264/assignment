package customers

import "assignment/invoices"

type Customer struct {
	Id        uint               `gorm:"column:id; primaryKey; autoIncrement"`
	FirstName string             `gorm:"column:first_name; not null"`
	LastName  string             `gorm:"column:last_name; not null"`
	Address   string             `gorm:"column:address; not null"`
	Phone     string             `gorm:"column:phone; not null"`
	Balance   float64            `gorm:"column:balance; not null"`
	Invoice   []invoices.Invoice `gorm:"foreignKey:Customer;references:Id"`
}
