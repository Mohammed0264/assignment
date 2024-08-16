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
	if result.Error != nil {
		return result.Error, 0
	}
	return nil, count
}
func (p *CustomerRepository) Delete(id uint) (error, int) {
	result := p.Db.Model(&Customer{}).Delete(&Customer{}, id)
	if result.Error != nil {
		return result.Error, 0
	}
	if result.RowsAffected > 0 {
		return nil, 1
	}
	return nil, 0

}
func (p *CustomerRepository) Find(id uint) (Customer, error, int64) {
	var foundCustomer Customer
	var count int64
	result := p.Db.Model(&Customer{}).Where("id = ?", &id).Find(&foundCustomer).Count(&count)
	if result.Error != nil {
		return Customer{}, result.Error, 0
	}
	if count == 1 {
		return foundCustomer, nil, 1

	}
	return foundCustomer, nil, 0
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
	if result.Error != nil {
		return result.Error, 0
	}
	if count == 0 {
		return nil, 0
	}
	return nil, 1
}

func (p *CustomerRepository) AddBalance(id uint, balance float64) (error, int64) {
	var count int64
	find, err, counter := p.Find(id)
	if err != nil {
		return err, 0
	}
	if counter == 0 {
		return nil, 0
	}
	newBalance := find.Balance + balance
	result := p.Db.Model(&Customer{}).Where("id=?", &id).Update("Balance", &newBalance).Count(&count)
	if result.Error != nil {
		return result.Error, 0
	}
	if count == 0 {
		return nil, 0
	}
	return nil, 1
}

func (p *CustomerRepository) SubtractBalance(id uint, cost float64) (error, int64) {
	var count int64
	find, err, counter := p.Find(id)
	if err != nil {
		return err, 0
	}
	if counter == 0 {
		return nil, 0
	}
	newBalance := find.Balance - cost
	if newBalance < 0 {
		return nil, 2
	}
	result := p.Db.Model(&Customer{}).Where("id=?", &id).Update("Balance", &newBalance).Count(&count)
	if result.Error != nil {
		return result.Error, 0
	}
	if count == 0 {
		return nil, 0
	}
	return nil, 1
}
