package suppliers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/microcosm-cc/bluemonday"
	"net/http"
)

type SupplierAPI struct {
	SupplierService SupplierService
}

func ProvideSupplierApi(p SupplierService) SupplierAPI {
	return SupplierAPI{SupplierService: p}
}
func (p *SupplierAPI) Create(c *gin.Context) {
	var supplierDto SupplierDto
	err := c.Bind(&supplierDto)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	supplierDto = sanitizeSupplier(supplierDto)
	err = validateSupplier(supplierDto)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = p.SupplierService.Create(ToSupplier(supplierDto))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusCreated, gin.H{"supplier": supplierDto})
}

func (p *SupplierAPI) Update(c *gin.Context) {

	var supplierDto SupplierDto
	err := c.Bind(&supplierDto)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	supplierDto = sanitizeSupplier(supplierDto)
	err = validateSupplier(supplierDto)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err, count := p.SupplierService.Update(ToSupplier(supplierDto))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if count == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"could not find supplier with id:": fmt.Sprint(supplierDto.Id)})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"supplier": supplierDto})
}
func (p *SupplierAPI) FindAll(c *gin.Context) {
	suppliers, err := p.SupplierService.FindAll()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"suppliers": ToSupplierDTOS(suppliers)})
}
func (p *SupplierAPI) FindByName(c *gin.Context) {
	name := c.Param("name")
	name = sanitizeName(name)
	suppliers, err := p.SupplierService.FindByName(name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"suppliers": ToSupplierDTOS(suppliers)})
}
func (p *SupplierAPI) Delete(c *gin.Context) {
	var supplierDto SupplierDto
	err := c.Bind(&supplierDto)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err, count := p.SupplierService.Delete(supplierDto.Id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if count == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"could not delete supplier because could not find supplier with id:": supplierDto.Id})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"supplier deleted with id": fmt.Sprint(supplierDto.Id)})
}
func validateSupplier(supplierDto SupplierDto) error {
	validate := validator.New()
	err := validate.Struct(supplierDto)
	if err != nil {
		return err
	}
	return nil
}
func sanitizeSupplier(supplierDto SupplierDto) SupplierDto {
	sanitize := bluemonday.StrictPolicy()
	supplierDto.Name = sanitize.Sanitize(supplierDto.Name)
	supplierDto.Phone = sanitize.Sanitize(supplierDto.Phone)
	return supplierDto
}
func sanitizeName(name string) string {
	sanitize := bluemonday.StrictPolicy()
	name = sanitize.Sanitize(name)
	return name
}
