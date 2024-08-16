package users

import "gorm.io/gorm"

type UserRepository struct {
	Db *gorm.DB
}

func ProvideUserRepository(db *gorm.DB) UserRepository {
	return UserRepository{Db: db}
}
func (p *UserRepository) Create(user User) error {
	return p.Db.Model(&User{}).Create(&user).Error
}
func (p *UserRepository) UpdateUserName(user User) (error, int64) {
	var count int64
	result := p.Db.Model(&User{}).Where("Id=?", &user.Id).Update("UserName", user.UserName).Count(&count)
	return result.Error, count
}
func (p *UserRepository) UpdatePassword(user User) (error, int64) {
	var count int64
	result := p.Db.Model(&User{}).Where("Id=?", &user.Id).Update("Password", user.Password).Count(&count)
	return result.Error, count
}
func (p *UserRepository) FindByUserName(userName string) []User {
	var users []User
	p.Db.Model(&User{}).Where("UserName=?", &userName).Find(&users)
	return users
}
func (p *UserRepository) FindAll() []User {
	var users []User
	p.Db.Model(&User{}).Find(&users)
	return users
}
func (p *UserRepository) Delete(id uint) (error, int64) {
	result := p.Db.Model(&User{}).Where("id=?", &id).Delete(&User{})
	return result.Error, result.RowsAffected
}
