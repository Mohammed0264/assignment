package users

func ToUser(userDto UserDto) User {
	return User{UserName: userDto.UserName, Password: userDto.Password}
}
func ToUserDto(user User) UserDto {
	return UserDto{Id: user.Id, UserName: user.UserName}
}
func ToUserDTOs(user []User) []UserDto {
	userDTOs := make([]UserDto, len(user))
	for index, value := range user {
		userDTOs[index] = ToUserDto(value)
	}
	return userDTOs
}
