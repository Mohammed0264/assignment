package users

import (
	"assignment/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type UserApi struct {
	UserService UserService
}

func ProvideUserApi(p UserService) UserApi {
	return UserApi{UserService: p}
}
func (p *UserApi) Create(c *gin.Context) {
	var userDto UserDto
	err := c.Bind(&userDto)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if userDto.UserName == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "username can't be empty"})
		return
	}
	if userDto.Password == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "password can't be empty"})
		return
	}
	userDto = sanitizeInput(userDto)
	err = validateInput(userDto)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(userDto.Password), bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "password": "could not hash password"})
		return
	}
	userDto.Password = string(hashPassword)
	err = p.UserService.Create(ToUser(userDto))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "user": "could not create new user"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"user": "created"})
}
func (p *UserApi) UpdateUserName(c *gin.Context) {
	var userDto UserDto
	err := c.Bind(&userDto)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userDto = sanitizeInput(userDto)
	/*err=validateInput(user)
	if err!=nil{
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}*/
	if userDto.UserName == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "username can't be empty"})
		return
	}
	err, count := p.UserService.UpdateUserName(ToUser(userDto))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "user": "could not update userName"})
		return
	}
	if count == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"user": "could not update userName"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"user": "updated"})
}
func (p *UserApi) UpdatePassword(c *gin.Context) {
	var userDto UserDto
	err := c.Bind(&userDto)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userDto = sanitizeInput(userDto)
	/*err=validateInput(user)
	if err!=nil{
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}*/
	if userDto.Password == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "password can't be empty"})
		return
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(userDto.Password), bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "password": "could not hash password"})
		return
	}
	userDto.Password = string(hashPassword)
	err, count := p.UserService.UpdatePassword(ToUser(userDto))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "user": "could not update password"})
		return
	}
	if count == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"user": "could not update password"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"user": "updated"})
}
func (p *UserApi) FindByUserName(c *gin.Context) {
	var userDto UserDto
	userName := c.Param("userName")
	userDto.UserName = userName
	sanitizeInput(userDto)
	var users []User
	users = p.UserService.FindByUserName(userDto.UserName)
	c.IndentedJSON(http.StatusOK, gin.H{"users": ToUserDTOs(users)})
}
func (p *UserApi) FindAll(c *gin.Context) {
	users := p.UserService.FindAll()
	c.IndentedJSON(http.StatusOK, gin.H{"users": ToUserDTOs(users)})
}
func (p *UserApi) Delete(c *gin.Context) {
	var userDto UserDto
	err := c.Bind(&userDto)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err, count := p.UserService.Delete(userDto.Id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if count == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"user": "could not delete user with id: " + fmt.Sprint(userDto.Id)})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"user": "deleted"})
}
func (p *UserApi) Login(c *gin.Context) {
	var userDto UserDto
	err := c.Bind(&userDto)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userDto = sanitizeInput(userDto)
	err = validateInput(userDto)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := p.UserService.FindUserNameLogin(userDto.UserName)
	if user.Id == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "wrong userName"})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userDto.Password))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "wrong password"})
		return
	}
	var userInformation middleware.UserInformation
	userInformation.Id = user.Id
	userInformation.UserName = user.UserName
	userInformation.Role = user.Role

	accessToken, err := middleware.CreateToken(userInformation)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "can not login contact admin"})
		return
	}

	refreshToken, err := middleware.GenerateRefreshToken(userInformation)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "can not login contact admin"})
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("refreshToken", refreshToken, 3600*24*7, "", "", false, true)
	c.IndentedJSON(http.StatusOK, gin.H{"token": accessToken})
}
func validateInput(user UserDto) error {
	validate := validator.New()
	return validate.Struct(&user)
}
func sanitizeInput(user UserDto) UserDto {
	sanitize := bluemonday.UGCPolicy()
	user.UserName = sanitize.Sanitize(user.UserName)
	user.Password = sanitize.Sanitize(user.Password)
	return user
}
