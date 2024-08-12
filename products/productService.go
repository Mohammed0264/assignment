package products

type ProductService struct {
	ProductRepository ProductRepository
}

func ProvideProductService(p ProductRepository) ProductService {
	return ProductService{ProductRepository: p}
}
func (p *ProductService) Create(product Product) error {
	err := p.ProductRepository.Create(product)
	if err != nil {
		return err
	} else {
		return nil
	}
}
func (p *ProductService) Update(product Product) (error, int64) {
	err, count := p.ProductRepository.Update(product)
	if err != nil {
		return err, 0
	} else {
		return nil, count
	}
}
func (p *ProductService) FindAll() []Product {
	return p.ProductRepository.FindAll()
}
func (p *ProductService) FindByName(name string) []Product {
	return p.ProductRepository.FindByName(name)
}
func (p *ProductService) Delete(id uint) (error, int64) {
	err, rowsAffected := p.ProductRepository.Delete(id)
	if err != nil {
		return err, 0
	} else {
		return nil, rowsAffected
	}
}
