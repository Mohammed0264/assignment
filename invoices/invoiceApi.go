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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = inputValidation(invoiceReceiver)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	invoiceReceiver = inputSanitization(invoiceReceiver)
	fmt.Println(invoiceReceiver)
	err = p.InvoiceService.Create(ToInvoice(invoiceReceiver), invoiceReceiver.InvoiceLine)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"invoice": "created"})
}
func (p *InvoiceApi) FindAll(c *gin.Context) {
	invoices := p.InvoiceService.FindAll()
	c.IndentedJSON(http.StatusOK, gin.H{"invoices": invoices})
}

func inputValidation(invoiceReceiver InvoiceReceiver) error {
	var validator = validator2.New()
	err := validator.Struct(invoiceReceiver)
	if err != nil {
		return err
	}
	return nil
}
func inputSanitization(invoiceReceiver InvoiceReceiver) InvoiceReceiver {
	var sanitize = bluemonday.UGCPolicy()
	invoiceReceiver.InvoiceUniqueId = sanitize.Sanitize(invoiceReceiver.InvoiceUniqueId)

	return invoiceReceiver
}
