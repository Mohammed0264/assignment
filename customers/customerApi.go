package customers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type CustomerApi struct {
	CustomerService CustomerService
}

func ProvideCustomerApi(p CustomerService) CustomerApi {
	return CustomerApi{CustomerService: p}
}
func (p *CustomerApi) Create(c *gin.Context) {
	var customerDto CustomerDto
	err := c.Bind(&customerDto)
	if err != nil {
		fmt.Println(err.Error())
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = p.CustomerService.Create(ToCustomer(customerDto))
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.IndentedJSON(http.StatusCreated, gin.H{"customer": customerDto})
}
func (p *CustomerApi) Update(c *gin.Context) {
	var customerDto CustomerDto
	err := c.Bind(&customerDto)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	err, count := p.CustomerService.Update(ToCustomer(customerDto))
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if count == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "customer not found with id:" + fmt.Sprint(customerDto.Id)})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"customer": customerDto})
}
func (p *CustomerApi) Delete(c *gin.Context) {
	var customerDto CustomerDto

	err := c.BindJSON(&customerDto)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err, rows := p.CustomerService.Delete(customerDto.Id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	if rows == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "There is No customer with id:" + fmt.Sprint(customerDto.Id)})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"success": "customer with id:" + fmt.Sprintf("%v", customerDto.Id) + " deleted"})
}
func (p *CustomerApi) FindById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	customer, err, count := p.CustomerService.GetCustomerById(uint(id))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	if count == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Customer with id:" + fmt.Sprint(id) + " not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"customer": ToCustomerDTO(customer)})
}
func (p *CustomerApi) FindAll(c *gin.Context) {
	customers, err := p.CustomerService.GetCustomers()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"customers": ToCustomerDTOs(customers)})
}
func (p *CustomerApi) UpdateBalance(c *gin.Context) {
	var customerDto CustomerDto
	err := c.Bind(&customerDto)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err, count := p.CustomerService.UpdateBalance(customerDto.Id, customerDto.Balance)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if count == 0 {
		c.IndentedJSON(http.StatusOK, gin.H{"failed": "customer with id:" + fmt.Sprintf("%v", customerDto.Id) + " does not exist"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"success": "balance updated"})
}
func (p *CustomerApi) SubtractBalance(c *gin.Context) {
	var customerDto CustomerDto
	err := c.Bind(&customerDto)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err, count := p.CustomerService.SubtractBalance(customerDto.Id, customerDto.Balance)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if count == 0 {
		c.IndentedJSON(http.StatusOK, gin.H{"failed": "customer with id:" + fmt.Sprintf("%v", customerDto.Id) + " does not exist"})
		return
	}
	if count == 2 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "you do not have enough balance"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"success": "balance updated"})
}
func (p *CustomerApi) AddBalance(c *gin.Context) {
	var customerDto CustomerDto
	err := c.Bind(&customerDto)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err, count := p.CustomerService.AddBalance(customerDto.Id, customerDto.Balance)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if count == 0 {
		c.IndentedJSON(http.StatusOK, gin.H{"failed": "customer with id:" + fmt.Sprintf("%v", customerDto.Id) + " does not exist"})
		return
	}
	if count == 2 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "you do not have enough balance"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"success": "balance updated"})
}
