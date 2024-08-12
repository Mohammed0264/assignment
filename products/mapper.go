package products

func ToProduct(productDto ProductDto) Product {
	return Product{Id: productDto.Id, Name: productDto.Name, Barcode: productDto.Barcode, QuantityOnHand: productDto.QuantityOnHand,
		Price: productDto.Price, Supplier: productDto.Supplier, ProductImage: productDto.ProductImage}

}
func ToProductDTO(product Product) ProductDto {
	return ProductDto{Id: product.Id, Name: product.Name, Barcode: product.Barcode, QuantityOnHand: product.QuantityOnHand,
		Price: product.Price, Supplier: product.Supplier, ProductImage: product.ProductImage}
}
func ToProductDTOs(products []Product) []ProductDto {
	productDTOs := make([]ProductDto, len(products))
	for index, value := range products {
		productDTOs[index] = ToProductDTO(value)
	}
	return productDTOs
}
