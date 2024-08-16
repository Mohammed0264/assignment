package suppliers

import (
	"errors"
	"strings"
)

type SupplierService struct {
	SupplierRepository SupplierRepository
}

func ProvideSupplierService(p SupplierRepository) SupplierService {
	return SupplierService{SupplierRepository: p}
}
func (p *SupplierService) Create(supplier Supplier) error {
	if strings.TrimSpace(supplier.Name) == "" {
		return errors.New("name should not be empty")
	}
	if strings.TrimSpace(supplier.Phone) == "" {
		return errors.New("phone should not be empty")

	}
	return p.SupplierRepository.Create(supplier)
}
func (p *SupplierService) Update(supplier Supplier) (error, int64) {

	return p.SupplierRepository.Update(supplier)
}
func (p *SupplierService) FindAll() ([]Supplier, error) {
	return p.SupplierRepository.FindAll()
}
func (p *SupplierService) FindByName(name string) ([]Supplier, error) {
	return p.SupplierRepository.FindByName(name)
}
func (p *SupplierService) Delete(id uint) (error, int) {
	return p.SupplierRepository.Delete(id)

}
