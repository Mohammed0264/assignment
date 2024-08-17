package customers

import "gorm.io/gorm"

type Customer struct {
	Id             uint    `gorm:"column:id; primaryKey; autoIncrement"`
	FirstName      string  `gorm:"column:first_name; not null"`
	LastName       string  `gorm:"column:last_name; not null"`
	Address        string  `gorm:"column:address; not null"`
	Phone          string  `gorm:"column:phone; not null"`
	Balance        float64 `gorm:"column:balance; not null"`
	gorm.DeletedAt `gorm:"column:deleted_at; index"`
	//	Invoice   []invoices.Invoice `gorm:"foreignKey:Customer;references:Id"`
}
