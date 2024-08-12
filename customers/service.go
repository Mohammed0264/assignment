package customers

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
func (p *CustomerService) Delete(id uint) (error, int) {
	return p.CustomerRepository.Delete(id)
}
func (p *CustomerService) AddBalance(id uint, balance float64) (error, int64) {
	return p.CustomerRepository.AddBalance(id, balance)
}
func (p *CustomerService) SubtractBalance(id uint, balance float64) (error, int64) {
	return p.CustomerRepository.SubtractBalance(id, balance)
}
