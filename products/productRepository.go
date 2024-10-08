package products

import (
	"gorm.io/gorm"
)

type ProductRepository struct {
	Db *gorm.DB
}

func ProvideProductRepository(db *gorm.DB) ProductRepository {
	return ProductRepository{Db: db}
}
func (p *ProductRepository) Create(product Product) error {
	result := p.Db.Model(&Product{}).Create(&product)
	return result.Error
}
func (p *ProductRepository) Update(product Product) (error, int64) {
	var count int64
	result := p.Db.Model(&Product{}).Where("Id=?", &product.Id).Omit("ProductImage").Updates(&product).Count(&count)
	return result.Error, count
}
func (p *ProductRepository) FindAll() []Product {
	var products []Product
	p.Db.Model(&Product{}).Find(&products)
	return products
}
func (p *ProductRepository) FindByName(name string) []Product {
	var products []Product
	p.Db.Model(&Product{}).Where("Name LIKE ?", "%"+name+"%").Find(&products)
	return products
}
func (p *ProductRepository) Delete(id uint) (error, int64) {
	result := p.Db.Model(&Product{}).Where("id=?", &id).Delete(&Product{})
	return result.Error, result.RowsAffected
}
func (p *ProductRepository) UpdateImage(id uint, image string) (error, int64) {
	var count int64
	result := p.Db.Model(&Product{}).Where("id=?", &id).Update("ProductImage", image).Count(&count)
	return result.Error, count

}
func (p *ProductRepository) FindById(id uint) Product {
	var product Product
	p.Db.Model(&Product{}).Where("Id=?", &id).Find(&product)
	return product
}
