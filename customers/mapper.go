package customers

func ToCustomer(customerDto CustomerDto) Customer {
	return Customer{Id: customerDto.Id, FirstName: customerDto.FirstName, LastName: customerDto.LastName,
		Address: customerDto.Address, Phone: customerDto.Phone, Balance: customerDto.Balance}
}
func ToCustomerDTO(customer Customer) CustomerDto {
	return CustomerDto{Id: customer.Id, FirstName: customer.FirstName, LastName: customer.LastName,
		Address: customer.Address, Balance: customer.Balance}
}
func ToCustomerDTOs(customers []Customer) []CustomerDto {
	customerDTOs := make([]CustomerDto, len(customers))
	for index, value := range customers {
		customerDTOs[index] = ToCustomerDTO(value)
	}
	return customerDTOs
}
