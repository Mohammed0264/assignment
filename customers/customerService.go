package customers

import "math"

type CustomerService struct {
	CustomerRepository CustomerRepository
}

func ProvideCustomerService(p CustomerRepository) CustomerService {
	return CustomerService{CustomerRepository: p}
}
func (p *CustomerService) Create(customer Customer) error {
	return p.CustomerRepository.Save(customer)

}
func (p *CustomerService) GetCustomers() ([]Customer, error) {
	return p.CustomerRepository.FindAll()
}
func (p *CustomerService) GetCustomerById(id uint) (Customer, error, int64) {

	return p.CustomerRepository.Find(id)

}
func (p *CustomerService) Update(customer Customer) (error, int64) {
	return p.CustomerRepository.Update(customer)
}
func (p *CustomerService) UpdateBalance(id uint, balance float64) (error, int64) {
	return p.CustomerRepository.UpdateBalance(id, balance)
}
func (p *CustomerService) Delete(id uint) (error, int64) {
	return p.CustomerRepository.Delete(id)
}
func (p *CustomerService) AddBalance(id uint, balance float64) (error, int64) {
	find, err, counter := p.CustomerRepository.Find(id)
	if err != nil {
		return err, 0
	}
	if counter == 0 {
		return nil, 0
	}
	// only include two numbers after pint like .00
	newBalance := find.Balance + balance
	newBalance = math.Round(newBalance*100) / 100
	return p.CustomerRepository.AddBalance(id, newBalance)
}
func (p *CustomerService) SubtractBalance(id uint, cost float64) (error, int64) {
	find, err, counter := p.CustomerRepository.Find(id)
	if err != nil {
		return err, 0
	}
	if counter == 0 {
		return nil, 0
	}
	// only include two numbers after pint like .00
	newBalance := find.Balance - cost
	newBalance = math.Round(newBalance*100) / 100
	if newBalance < 0 {
		return nil, 2
	}
	return p.CustomerRepository.SubtractBalance(id, newBalance)
}
