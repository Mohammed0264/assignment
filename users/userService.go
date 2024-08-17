package users

type UserService struct {
	UserRepository UserRepository
}

func ProvideUserService(p UserRepository) UserService {
	return UserService{UserRepository: p}
}
func (p *UserService) Create(user User) error {
	return p.UserRepository.Create(user)
}
func (p *UserService) UpdateUserName(user User) (error, int64) {
	return p.UserRepository.UpdateUserName(user)
}
func (p *UserService) UpdatePassword(user User) (error, int64) {
	return p.UserRepository.UpdatePassword(user)
}
func (p *UserService) FindByUserName(userName string) []User {
	return p.UserRepository.FindByUserName(userName)
}
func (p *UserService) FindUserNameLogin(userName string) User {
	return p.UserRepository.FindUserNameLogin(userName)
}
func (p *UserService) FindAll() []User {
	return p.UserRepository.FindAll()
}
func (p *UserService) Delete(id uint) (error, int64) {
	return p.UserRepository.Delete(id)
}
