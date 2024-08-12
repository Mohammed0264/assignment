package products

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ProductApi struct {
	ProductService ProductService
}

func ProvideProductApi(p ProductService) ProductApi {
	return ProductApi{ProductService: p}
}
func (p *ProductApi) Create(c *gin.Context) {
	var productDto ProductDto
	err := c.Bind(&productDto)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = p.ProductService.Create(ToProduct(productDto))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.IndentedJSON(http.StatusCreated, gin.H{"product": productDto})
}

func (p *ProductApi) Update(c *gin.Context) {
	var productDto ProductDto
	err := c.Bind(&productDto)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	er, count := p.ProductService.Update(ToProduct(productDto))
	if er != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	} else if count == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "could not update product check item id"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"product": productDto})
}
func (p *ProductApi) FindAll(c *gin.Context) {
	products := p.ProductService.FindAll()
	c.IndentedJSON(http.StatusOK, gin.H{"products": ToProductDTOs(products)})
}
func (p *ProductApi) FindByName(c *gin.Context) {
	name := c.Param("name")
	products := p.ProductService.FindByName(name)
	c.IndentedJSON(http.StatusOK, gin.H{"products": ToProductDTOs(products)})
}
func (p *ProductApi) Delete(c *gin.Context) {
	var productDto ProductDto
	err := c.Bind(&productDto)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	er, affectedRows := p.ProductService.Delete(productDto.Id)
	if er != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	} else if affectedRows == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "could not delete product check item id"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"product": "product deleted with id: " + fmt.Sprint(productDto.Id)})
}
