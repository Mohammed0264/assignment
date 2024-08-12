package products

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
func (p *ProductApi) UpdateImage(c *gin.Context) {
	id, err := strconv.Atoi(c.PostForm("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	originalImage := c.PostForm("originalImage")
	if originalImage != "Null1" {
		_, err = os.Stat(originalImage)
		if err == nil {
			err = os.Remove(originalImage)
			fmt.Println("image removed")
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
		}
	}

	image, err := c.FormFile("image")
	fmt.Println("1")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	imageName := image.Filename
	extension := strings.ToLower(filepath.Ext(imageName))
	valid := checkFileExtension(extension)
	if !valid {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid file extension"})
		return
	}

	file := filepath.Join("./images/", imageName)
	counter := 0
	fmt.Println("2")

	for {
		_, err = os.Stat(file)
		if err == nil {

			file = filepath.Join("./images/", strconv.Itoa(counter)+imageName)
			counter++
			continue
		} else {
			break
		}
	}

	err, count := p.ProductService.UpdateImage(uint(id), file)
	fmt.Println("3")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
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

func checkFileExtension(extension string) bool {
	validExtension := []string{".jpg", ".png", ".jpeg"}
	for _, value := range validExtension {
		if extension == value {
			return true
		}
	}
	return false
}
