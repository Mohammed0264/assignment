package invoices

import (
	"fmt"
	"github.com/gin-gonic/gin"
	validator2 "github.com/go-playground/validator/v10"
	"github.com/microcosm-cc/bluemonday"
	"net/http"
)

type InvoiceApi struct {
	InvoiceService InvoiceService
}

func ProvideInvoiceApi(p InvoiceService) InvoiceApi {
	return InvoiceApi{InvoiceService: p}
}
func (p *InvoiceApi) Create(c *gin.Context) {
	var invoiceReceiver InvoiceReceiver
	err := c.Bind(&invoiceReceiver)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	invoiceReceiver = inputSanitizationForCreate(invoiceReceiver)
	err = inputValidationForCreate(invoiceReceiver)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(invoiceReceiver)
	err = p.InvoiceService.Create(ToInvoice(invoiceReceiver), invoiceReceiver.InvoiceLine)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"invoice": "created"})
}
func (p *InvoiceApi) FindAll(c *gin.Context) {
	invoices := p.InvoiceService.FindAll()
	c.IndentedJSON(http.StatusOK, gin.H{"invoices": invoices})
}
func (p *InvoiceApi) Update(c *gin.Context) {
	var updateInvoice InvoiceUpdate
	err := c.Bind(&updateInvoice)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updateInvoice = inputSanitizationForUpdate(updateInvoice)

	err = inputValidationForUpdate(updateInvoice)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err, count := p.InvoiceService.Update(updateInvoice)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, gin.H{"invoice": count})
}
func (p *InvoiceApi) Delete(c *gin.Context) {
	var invoiceDto InvoiceDto
	err := c.Bind(&invoiceDto)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err, count := p.InvoiceService.Delete(invoiceDto.Id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if count == 0 {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "can not delete invoice"})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"invoice": "invoice deleted"})

}
func inputValidationForCreate(invoiceReceiver InvoiceReceiver) error {
	var validator = validator2.New()
	err := validator.Struct(&invoiceReceiver)
	if err != nil {
		return err
	}
	return nil
}
func inputValidationForUpdate(invoiceUpdate InvoiceUpdate) error {
	var validator = validator2.New()
	return validator.Struct(&invoiceUpdate)
}
func inputSanitizationForCreate(invoiceReceiver InvoiceReceiver) InvoiceReceiver {
	var sanitize = bluemonday.UGCPolicy()
	invoiceReceiver.InvoiceUniqueId = sanitize.Sanitize(invoiceReceiver.InvoiceUniqueId)
	for index, _ := range invoiceReceiver.InvoiceLine {
		invoiceReceiver.InvoiceLine[index].ItemName = sanitize.Sanitize(invoiceReceiver.InvoiceLine[index].ItemName)
	}
	return invoiceReceiver
}
func inputSanitizationForUpdate(invoiceUpdate InvoiceUpdate) InvoiceUpdate {
	var sanitize = bluemonday.UGCPolicy()
	invoiceUpdate.InvoiceUniqueId = sanitize.Sanitize(invoiceUpdate.InvoiceUniqueId)
	for index, _ := range invoiceUpdate.InvoiceLineDto {
		invoiceUpdate.InvoiceLineDto[index].ItemName = sanitize.Sanitize(invoiceUpdate.InvoiceLineDto[index].ItemName)
	}
	for index, _ := range invoiceUpdate.UpdateInvoiceLine {
		invoiceUpdate.UpdateInvoiceLine[index].ItemName = sanitize.Sanitize(invoiceUpdate.UpdateInvoiceLine[index].ItemName)
	}
	return invoiceUpdate
}
