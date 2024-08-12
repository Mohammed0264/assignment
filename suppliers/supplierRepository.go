package suppliers

import "gorm.io/gorm"

type SupplierRepository struct {
	Db *gorm.DB
}

func ProvideSupplierRepository(db *gorm.DB) SupplierRepository {
	return SupplierRepository{Db: db}
}
func (p *SupplierRepository) Create(supplier Supplier) error {
	result := p.Db.Model(&Supplier{}).Create(&supplier)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (p *SupplierRepository) Update(supplier Supplier) (error, int64) {
	var count int64
	result := p.Db.Model(&Supplier{}).Where("Id=?", &supplier.Id).Updates(&supplier).Count(&count)
	if result.Error != nil {
		return result.Error, 0
	}
	return nil, count
}
func (p *SupplierRepository) Delete(id uint) (error, int) {
	result := p.Db.Model(&Supplier{}).Where("id=?", &id).Delete(&Supplier{})
	if result.Error != nil {
		return result.Error, 0
	}
	if result.RowsAffected > 0 {
		return nil, 1
	}
	return nil, 0
}
func (p *SupplierRepository) FindAll() ([]Supplier, error) {
	var suppliers []Supplier
	result := p.Db.Model(&Supplier{}).Find(&suppliers)
	if result.Error != nil {
		return nil, result.Error
	}
	return suppliers, nil
}
func (p *SupplierRepository) FindByName(name string) ([]Supplier, error) {
	var suppliers []Supplier
	result := p.Db.Model(&Supplier{}).Where("Name LIKE ?", "%"+name+"%").Find(&suppliers)
	if result.Error != nil {
		return nil, result.Error
	}
	return suppliers, nil
}
