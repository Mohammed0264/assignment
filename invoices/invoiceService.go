package invoices

import (
	"assignment/customers"
	"assignment/invoiceLines"
	"assignment/products"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var InitProductApiReceiver products.ProductApi
var InitCustomerApiReceiver customers.CustomerApi
var InitInvoiceLineServiceReceiver invoiceLines.InvoiceLineService

type InvoiceService struct {
	InvoiceRepository InvoiceRepository
}

func ProvideInvoiceService(p InvoiceRepository) InvoiceService {
	return InvoiceService{InvoiceRepository: p}
}
func (p *InvoiceService) Create(invoice Invoice, invoiceLine []invoiceLines.InvoiceLineDto) error {
	customerService := InitCustomerApiReceiver.CustomerService

	customer, err, count := customerService.GetCustomerById(invoice.Customer)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("customer not found")
	}
	productService := InitProductApiReceiver.ProductService
	var product = make([]products.Product, len(invoiceLine))
	var totalPrice float64
	for index, value := range invoiceLine {
		product[index] = productService.FindById(value.ItemId)
		if product[index].QuantityOnHand-invoiceLine[index].Quantity < 0 {
			return errors.New(product[index].Name + ": quantity is negative")
		}
		invoiceLine[index].LinePrice = product[index].Price
		totalPrice = totalPrice + (product[index].Price * invoiceLine[index].Quantity)
	}
	if customer.Balance-totalPrice < 0 {
		return errors.New(customer.FirstName + " " + customer.LastName + " does not have enough balance")
	}
	retrieveLastInvoice := p.InvoiceRepository.FindLastInvoice()

	if retrieveLastInvoice.Id == 0 {
		retrieveLastInvoice.InvoiceUniqueId = "0000-0000"
	}

	time1 := time.Now()
	year := time1.Year()

	lastYear, err := strconv.Atoi(strings.Split(retrieveLastInvoice.InvoiceUniqueId, "-")[0])
	if err != nil {
		return err
	}

	uniqueId, err := strconv.Atoi(strings.Split(retrieveLastInvoice.InvoiceUniqueId, "-")[1])
	if err != nil {
		return err
	}

	var newId string
	if year > lastYear {
		uniqueId = 1
		formatId := fmt.Sprintf("%04d", uniqueId)
		newId = strconv.Itoa(year) + "-" + formatId
	} else {
		uniqueId = uniqueId + 1
		if uniqueId <= 1000 {
			formatId := fmt.Sprintf("%04d", uniqueId)
			newId = strconv.Itoa(lastYear) + "-" + formatId
		} else {
			formatId := fmt.Sprint(uniqueId)
			newId = strconv.Itoa(lastYear) + "-" + formatId
		}

	}
	fmt.Println(newId)
	invoice.InvoiceUniqueId = newId
	invoice.InvoiceTotal = totalPrice
	err, lastInvoice := p.InvoiceRepository.Create(invoice)
	if err != nil {
		return err
	}

	for index, value := range invoiceLine {
		invoiceLine[index].InvoiceId = lastInvoice.Id

		err = InitInvoiceLineServiceReceiver.Create(invoiceLine[index])
		if err != nil {
			err, _ = p.Delete(newId)
			if err != nil {
				return errors.New("could not complete invoice request and could not delete please contact admin " +
					"as fast as possible")
			}
			return err
		}
		product[index].QuantityOnHand = product[index].QuantityOnHand - value.Quantity
		err, _ = productService.Update(product[index])
		if err != nil {
			return err
		}

		err, _ = InitCustomerApiReceiver.CustomerService.SubtractBalance(customer.Id, totalPrice)
		if err != nil {
			return err
		}
	}
	return nil
}
func (p *InvoiceService) Update(invoice Invoice) (error, int64) {
	return p.InvoiceRepository.Update(invoice)
}
func (p *InvoiceService) FindAll() []InvoiceSender {
	invoices := p.InvoiceRepository.FindAll()
	invoiceSender := make([]InvoiceSender, len(invoices))
	fmt.Println(invoices)

	invoiceLineService := InitInvoiceLineServiceReceiver
	customerService := InitCustomerApiReceiver
	productService := InitProductApiReceiver
	fmt.Println(invoiceLineService)
	fmt.Println(productService)
	for index, _ := range invoices {
		invoiceLine := invoiceLines.ToInvoiceDTOs(invoiceLineService.FindByInvoiceId(invoices[index].Id))
		fmt.Println(invoiceLine)
		customer, _, _ := customerService.CustomerService.GetCustomerById(invoices[index].Customer)
		fmt.Println(customer)
		invoiceSender[index].Id = invoices[index].Id
		invoiceSender[index].InvoiceUniqueId = invoices[index].InvoiceUniqueId
		invoiceSender[index].Customer = invoices[index].Customer
		invoiceSender[index].CustomerName = customer.FirstName + " " + customer.LastName
		invoiceSender[index].InvoiceDate = invoices[index].InvoiceDate
		invoiceSender[index].InvoiceTotal = invoices[index].InvoiceTotal
		for index2, _ := range invoiceLine {
			product := productService.ProductService.FindById(invoiceLine[index2].ItemId)
			invoiceLine[index2].ItemName = product.Name
		}
		invoiceSender[index].InvoiceLineDto = invoiceLine

	}

	return invoiceSender
}
func (p *InvoiceService) FindById(invoiceUniqueId string) (Invoice, error) {
	return p.InvoiceRepository.FindById(invoiceUniqueId)
}
func (p *InvoiceService) Delete(invoiceUniqueId string) (error, int64) {
	return p.InvoiceRepository.Delete(invoiceUniqueId)
}
