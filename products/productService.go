package products

import (
	"errors"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type ProductService struct {
	ProductRepository ProductRepository
}

func ProvideProductService(p ProductRepository) ProductService {
	return ProductService{ProductRepository: p}
}
func (p *ProductService) Create(product Product) error {
	err := p.ProductRepository.Create(product)
	if err != nil {
		return err
	} else {
		return nil
	}
}
func (p *ProductService) Update(product Product) (error, int64) {
	err, count := p.ProductRepository.Update(product)
	if err != nil {
		return err, 0
	} else {
		return nil, count
	}
}
func (p *ProductService) FindAll() []Product {
	return p.ProductRepository.FindAll()
}
func (p *ProductService) FindByName(name string) []Product {
	return p.ProductRepository.FindByName(name)
}
func (p *ProductService) Delete(id uint) (error, int64) {
	err, rowsAffected := p.ProductRepository.Delete(id)
	if err != nil {
		return err, 0
	} else {
		return nil, rowsAffected
	}
}
func (p *ProductService) UpdateImage(id uint, image *multipart.FileHeader, originalImage string) (error, int64, string) {
	imageName := image.Filename
	extension := strings.ToLower(filepath.Ext(imageName))
	valid := checkFileExtension(extension)
	if !valid {
		return errors.New("invalid file extension"), 0, ""
	}

	file := filepath.Join("./images/", imageName)
	counter := 0

	for {
		_, err := os.Stat(file)
		if err == nil {

			file = filepath.Join("./images/", strconv.Itoa(counter)+imageName)
			counter++
			continue
		} else {
			break
		}
	}

	err, count := p.ProductRepository.UpdateImage(id, file)
	if err != nil {
		return err, 0, ""
	} else {
		return nil, count, file
	}
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
func removeOriginalImage(image string) error {
	if image != "Null1" {
		_, err := os.Stat(image)
		if err == nil {
			err = os.Remove(image)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
