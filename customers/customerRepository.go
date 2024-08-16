package customers

import (
	"gorm.io/gorm"
)

type CustomerRepository struct {
	Db *gorm.DB
}

func ProvideCustomerRepository(db *gorm.DB) CustomerRepository {
	return CustomerRepository{Db: db}
}
func (p *CustomerRepository) Save(customer Customer) error {
	result := p.Db.Model(&Customer{}).Create(&customer)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (p *CustomerRepository) Update(customer Customer) (error, int64) {
	var count int64
	result := p.Db.Model(&Customer{}).Where("Id=?", customer.Id).Omit("Balance", "Id").Updates(&customer).Count(&count)
	return result.Error, count
}
func (p *CustomerRepository) Delete(id uint) (error, int64) {
	result := p.Db.Model(&Customer{}).Delete(&Customer{}, id)
	return result.Error, result.RowsAffected
}
func (p *CustomerRepository) Find(id uint) (Customer, error, int64) {
	var foundCustomer Customer
	var count int64
	result := p.Db.Model(&Customer{}).Where("id = ?", &id).Find(&foundCustomer).Count(&count)
	return foundCustomer, result.Error, count
}
func (p *CustomerRepository) FindAll() ([]Customer, error) {
	var customers []Customer
	result := p.Db.Model(&Customer{}).Find(&customers)
	if result.Error != nil {
		return nil, result.Error
	}
	return customers, nil
}
func (p *CustomerRepository) UpdateBalance(id uint, balance float64) (error, int64) {
	var count int64
	result := p.Db.Model(&Customer{}).Where("id=?", &id).Update("Balance", balance).Count(&count)
	return result.Error, count
}

func (p *CustomerRepository) AddBalance(id uint, newBalance float64) (error, int64) {
	var count int64
	result := p.Db.Model(&Customer{}).Where("id=?", &id).Update("Balance", &newBalance).Count(&count)
	return result.Error, count
}

func (p *CustomerRepository) SubtractBalance(id uint, newBalance float64) (error, int64) {
	var count int64
	result := p.Db.Model(&Customer{}).Where("id=?", &id).Update("Balance", &newBalance).Count(&count)
	return result.Error, count
}
