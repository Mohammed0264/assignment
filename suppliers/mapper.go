package suppliers

func ToSupplier(supplierDto SupplierDto) Supplier {
	return Supplier{Id: supplierDto.Id, Name: supplierDto.Name, Phone: supplierDto.Phone}
}
func ToSupplierDTO(supplier Supplier) SupplierDto {
	return SupplierDto{Id: supplier.Id, Name: supplier.Name, Phone: supplier.Phone}
}
func ToSupplierDTOS(supplier []Supplier) []SupplierDto {
	supplierDTOs := make([]SupplierDto, len(supplier))
	for index, value := range supplier {
		supplierDTOs[index] = ToSupplierDTO(value)
	}
	return supplierDTOs
}
