package products

import (
	"fmt"
	"github.com/gin-gonic/gin"
	validator2 "github.com/go-playground/validator/v10"
	"github.com/microcosm-cc/bluemonday"
	"net/http"
	"strconv"
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
	// sanitize inputs to remove scripts because of xss and sql injection
	productDto = sanitizeInput(productDto)
	// check input validation
	err = validateInput(productDto)
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
	// sanitize inputs to remove scripts because of xss and sql injection
	productDto = sanitizeInput(productDto)
	// check input validation
	err = validateInput(productDto)
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
	var validator = validator2.New()
	var sanitize = bluemonday.UGCPolicy()
	err := validator.Var(name, "required")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	name = sanitize.Sanitize(name)
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

func (p *ProductApi) UpdateImage(c *gin.Context) {
	var validator = validator2.New()
	var sanitize = bluemonday.UGCPolicy()
	id, err := strconv.Atoi(c.PostForm("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	originalImage := c.PostForm("originalImage")
	sanitize.Sanitize(originalImage)

	err = validator.Var(originalImage, "required")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	image, err := c.FormFile("image")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err, count, file := p.ProductService.UpdateImage(uint(id), image, originalImage)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if count == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "could not update product image"})
		return
	}

	err = c.SaveUploadedFile(image, file)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"image": "image uploaded"})

}

func sanitizeInput(productDto ProductDto) ProductDto {
	san := bluemonday.UGCPolicy()
	productDto.Name = san.Sanitize(productDto.Name)
	productDto.Barcode = san.Sanitize(productDto.Barcode)
	productDto.ProductImage = san.Sanitize(productDto.ProductImage)
	return productDto
}
func validateInput(productDto ProductDto) error {
	var validator = validator2.New()
	err := validator.Struct(productDto)
	if err != nil {
		return err
	}
	return nil
}
