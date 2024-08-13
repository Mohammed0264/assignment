package invoices

import (
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
	var invoiceDto InvoiceDto
	err := c.Bind(&invoiceDto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = inputValidation(invoiceDto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	invoiceDto = inputSanitization(invoiceDto)
	err = p.InvoiceService.Create(ToInvoice(invoiceDto))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"invoice": invoiceDto})
}

func inputValidation(invoiceDto InvoiceDto) error {
	var validator = validator2.New()
	err := validator.Struct(invoiceDto)
	if err != nil {
		return err
	}
	return nil
}
func inputSanitization(invoiceDto InvoiceDto) InvoiceDto {
	var sanitize = bluemonday.UGCPolicy()
	invoiceDto.InvoiceUniqueId = sanitize.Sanitize(invoiceDto.InvoiceUniqueId)
	return invoiceDto
}
