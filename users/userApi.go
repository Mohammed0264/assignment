package users

import (
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
	var user User
	err := c.Bind(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if user.UserName == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "username can't be empty"})
		return
	}
	if user.Password == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "password can't be empty"})
		return
	}
	user = sanitizeInput(user)
	err = validateInput(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "password": "could not hash password"})
	}
	user.Password = string(hashPassword)
	err = p.UserService.Create(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "user": "could not create new user"})
	}
	c.JSON(http.StatusCreated, gin.H{"user": "created"})
}
func (p *UserApi) UpdateUserName(c *gin.Context) {
	var user User
	err := c.Bind(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user = sanitizeInput(user)
	/*err=validateInput(user)
	if err!=nil{
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}*/
	if user.UserName == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "username can't be empty"})
		return
	}
	err, count := p.UserService.UpdateUserName(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "user": "could not update userName"})
		return
	}
	if count == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"user": "could not update userName"})
	}
	c.IndentedJSON(http.StatusOK, gin.H{"user": "updated"})
}
func (p *UserApi) UpdatePassword(c *gin.Context) {
	var user User
	err := c.Bind(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user = sanitizeInput(user)
	/*err=validateInput(user)
	if err!=nil{
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}*/
	if user.Password == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "password can't be empty"})
		return
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "password": "could not hash password"})
	}
	user.Password = string(hashPassword)
	err, count := p.UserService.UpdatePassword(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error(), "user": "could not update password"})
		return
	}
	if count == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"user": "could not update password"})
	}
	c.IndentedJSON(http.StatusOK, gin.H{"user": "updated"})
}
func (p *UserApi) FindByUserName(c *gin.Context) {
	var user User
	userName := c.Param("userName")
	user.UserName = userName
	sanitizeInput(user)
	var users []User
	users = p.UserService.FindByUserName(user.UserName)
	c.IndentedJSON(http.StatusOK, gin.H{"users": ToUserDTOs(users)})
}
func (p *UserApi) FindAll(c *gin.Context) {
	var users []User
	users = p.UserService.FindAll()
	c.IndentedJSON(http.StatusOK, gin.H{"users": ToUserDTOs(users)})
}
func (p *UserApi) Delete(c *gin.Context) {
	var user User
	err := c.Bind(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err, count := p.UserService.Delete(user.Id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if count == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"user": "could not delete user with id: " + fmt.Sprint(user.Id)})
	}
	c.IndentedJSON(http.StatusOK, gin.H{"user": "deleted"})
}
func validateInput(user User) error {
	validate := validator.New()
	return validate.Struct(&user)
}
func sanitizeInput(user User) User {
	sanitize := bluemonday.UGCPolicy()
	user.UserName = sanitize.Sanitize(user.UserName)
	user.Password = sanitize.Sanitize(user.Password)
	user.Role = sanitize.Sanitize(user.Role)
	return user
}
