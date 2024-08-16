package invoices

import (
	"assignment/customers"
	"assignment/invoiceLines"
	"assignment/products"
	"errors"
	"fmt"
	"math"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var InitProductApiReceiver products.ProductApi
var InitCustomerApiReceiver customers.CustomerApi
var InitInvoiceLineServiceReceiver invoiceLines.InvoiceLineService
var updateCustomerId uint
var totalPriceUpdated float64

type InvoiceService struct {
	InvoiceRepository InvoiceRepository
}

func ProvideInvoiceService(p InvoiceRepository) InvoiceService {
	return InvoiceService{InvoiceRepository: p}
}
func (p *InvoiceService) Create(invoice Invoice, invoiceLineDto []invoiceLines.InvoiceLineDto) error {
	tax, err := strconv.ParseFloat(os.Getenv("TAX_RATE"), 64)
	if err != nil {
		return err
	}
	addTax := 1 + tax/100
	threshold, err := strconv.ParseFloat(os.Getenv("TAX_THRESHOLD"), 64)
	if err != nil {
		return err
	}
	customerService := InitCustomerApiReceiver.CustomerService

	customer, err, count := customerService.GetCustomerById(invoice.Customer)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("customer not found")
	}
	productService := InitProductApiReceiver.ProductService
	var product = make([]products.Product, len(invoiceLineDto))
	var totalPrice float64
	for index, value := range invoiceLineDto {
		product[index] = productService.FindById(value.ItemId)
		if product[index].QuantityOnHand-invoiceLineDto[index].Quantity < 0 {
			return errors.New(product[index].Name + ": quantity is negative")
		}
		invoiceLineDto[index].LinePrice = product[index].Price
		totalPrice = totalPrice + (product[index].Price * invoiceLineDto[index].Quantity)
	}
	if totalPrice > threshold {
		totalPrice = totalPrice * addTax
		totalPriceUpdated = math.Round(totalPriceUpdated*100) / 100
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
	invoice.InvoiceUniqueId = newId
	fmt.Println("last")
	fmt.Println(totalPriceUpdated)
	invoice.InvoiceTotal = totalPrice
	err, lastInvoice := p.InvoiceRepository.Create(invoice)
	if err != nil {
		return err
	}

	for index, value := range invoiceLineDto {
		invoiceLineDto[index].InvoiceId = lastInvoice.Id

		err = InitInvoiceLineServiceReceiver.Create(invoiceLineDto[index])
		if err != nil {
			err, _ = p.Delete(p.InvoiceRepository.FindLastInvoice().Id)
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

	}
	err, _ = InitCustomerApiReceiver.CustomerService.SubtractBalance(customer.Id, totalPrice)
	if err != nil {
		return err
	}
	return nil
}
func (p *InvoiceService) Update(invoiceUpdate InvoiceUpdate) (error, int64) {
	// here update
	err, count := p.updateAndDeleteInvoice(invoiceUpdate)
	if err != nil {
		return err, 0
	}
	if count == 4 {
		err, count = p.updateAndCreateInvoice(invoiceUpdate)
		if err != nil {
			return err, 0
		}
	}
	if count == 5 {
		err, _ = p.updateInvoiceEqualLength(invoiceUpdate)
		if err != nil {
			return err, 0
		}
	}

	return nil, 0
	//return p.InvoiceRepository.Update(invoice)
}
func (p *InvoiceService) FindAll() []InvoiceSender {
	invoices := p.InvoiceRepository.FindAll()
	invoiceSender := make([]InvoiceSender, len(invoices))

	invoiceLineService := InitInvoiceLineServiceReceiver
	customerService := InitCustomerApiReceiver
	productService := InitProductApiReceiver
	for index, _ := range invoices {
		invoiceLine := invoiceLines.ToInvoiceDTOs(invoiceLineService.FindByInvoiceId(invoices[index].Id))
		customer, _, _ := customerService.CustomerService.GetCustomerById(invoices[index].Customer)
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
func (p *InvoiceService) FindByInvoiceUniqueId(invoiceUniqueId string) (Invoice, error) {
	return p.InvoiceRepository.FindByInvoiceUniqueId(invoiceUniqueId)
}
func (p *InvoiceService) Delete(invoiceId uint) (error, int64) {
	invoiceLine := InitInvoiceLineServiceReceiver
	product := InitProductApiReceiver.ProductService
	customer := InitCustomerApiReceiver.CustomerService
	invoiceLinesReceived := invoiceLine.FindByInvoiceId(invoiceId)
	invoice, err := p.InvoiceRepository.FindById(invoiceId)
	if err != nil {
		return errors.New("could not find invoice"), 0
	}
	if invoice.Id == 0 {
		return errors.New("could not find invoice"), 0
	}
	customerId := invoice.Customer
	counter := 1
	for index, _ := range invoiceLinesReceived {
		updateProduct := product.FindById(invoiceLinesReceived[index].ItemId)
		updateProduct.QuantityOnHand += invoiceLinesReceived[index].Quantity
		err, count := product.Update(updateProduct)
		if err != nil {
			return err, 0
		}
		if count == 0 {
			return errors.New("error on " + fmt.Sprint(counter) + " try of update Product"), 0
		}
		err, count = invoiceLine.Delete(invoiceLinesReceived[index].Id)
		if err != nil {
			return err, 0
		}

		if count == 0 {
			return errors.New("error on " + fmt.Sprint(counter) + " try of delete invoice line"), 0
		}

		counter++
	}
	err, count := customer.AddBalance(customerId, invoice.InvoiceTotal)
	if err != nil {
		return err, 0
	}
	if count == 0 {
		return errors.New("error on " + fmt.Sprint(counter) + " try of add balance"), 0
	}
	err, count = p.InvoiceRepository.Delete(invoiceId)
	if err != nil {
		return err, 0
	}
	return nil, count
}

func (p *InvoiceService) updateAndDeleteInvoice(invoiceUpdate InvoiceUpdate) (error, int64) {
	tax, err := strconv.ParseFloat(os.Getenv("TAX_RATE"), 64)
	if err != nil {
		return err, 0
	}
	addTax := 1 + tax/100
	threshold, err := strconv.ParseFloat(os.Getenv("TAX_THRESHOLD"), 64)
	if err != nil {
		return err, 0
	}
	updateCustomerId = 0
	totalPriceUpdated = 0
	customerService := InitCustomerApiReceiver.CustomerService
	invoiceLineService := InitInvoiceLineServiceReceiver
	productService := InitProductApiReceiver.ProductService
	// check available quantities
	for index, _ := range invoiceUpdate.UpdateInvoiceLine {
		product := productService.FindById(invoiceUpdate.UpdateInvoiceLine[index].ItemId)
		if product.QuantityOnHand-invoiceUpdate.UpdateInvoiceLine[index].Quantity < 0 {
			return errors.New("we do not have enough quantity of item:" + invoiceUpdate.UpdateInvoiceLine[index].ItemName), 0
		}
	}

	//update total price
	for index, _ := range invoiceUpdate.UpdateInvoiceLine {
		item := productService.FindById(invoiceUpdate.UpdateInvoiceLine[index].ItemId)
		totalPriceUpdated += item.Price * invoiceUpdate.UpdateInvoiceLine[index].Quantity
	}
	if totalPriceUpdated >= threshold {
		totalPriceUpdated = totalPriceUpdated * addTax
		totalPriceUpdated = math.Round(totalPriceUpdated*100) / 100
	}
	// check for customer update
	if invoiceUpdate.UpdateCustomer != invoiceUpdate.OriginalCustomer {
		customer, err, count := customerService.GetCustomerById(invoiceUpdate.UpdateCustomer)
		if err != nil {
			return err, 0
		}
		if count == 0 {
			return errors.New("new customer not found"), 0
		}
		updateCustomerId = customer.Id

		if customer.Balance-totalPriceUpdated < 0 {
			return errors.New("updated customer has not enough balance"), 0
		}
	} else {
		updateCustomerId = invoiceUpdate.OriginalCustomer

		customer, err, count := customerService.GetCustomerById(invoiceUpdate.OriginalCustomer)
		customerBalance := invoiceUpdate.InvoiceTotal + customer.Balance
		if err != nil {
			return err, 0
		}
		if count == 0 {
			return errors.New("could not find data for old customer"), 0
		}
		if customerBalance-totalPriceUpdated < 0 {
			return errors.New("customer has not have enough money"), 0
		}

	}

	notifier := 0
	counter := 0
	// check for delete invoiceLines and create
	if len(invoiceUpdate.UpdateInvoiceLine) < len(invoiceUpdate.InvoiceLineDto) {
		var resizeInvoiceLineDto []invoiceLines.InvoiceLineDto
		// check for create
		for i := 0; i < len(invoiceUpdate.UpdateInvoiceLine); i++ {
			fmt.Println("insdie loop of create")
			if invoiceUpdate.UpdateInvoiceLine[i].Id == 0 {
				fmt.Println("inside create")
				err := invoiceLineService.Create(invoiceUpdate.UpdateInvoiceLine[i])
				if err != nil {
					return errors.New("can not create new invoiceLine1"), 0

				}
				product := productService.FindById(invoiceUpdate.UpdateInvoiceLine[i].ItemId)
				product.QuantityOnHand = product.QuantityOnHand - invoiceUpdate.UpdateInvoiceLine[i].Quantity
				err, count := productService.Update(product)
				if err != nil {
					return err, 0
				}
				if count == 0 {
					return errors.New("could not update product1"), 0
				}

				//	break

			}
		}

		// check for update and delete
		for b := 0; b < len(invoiceUpdate.InvoiceLineDto); b++ {
			for i := 0; i < len(invoiceUpdate.UpdateInvoiceLine); i++ {

				if invoiceUpdate.InvoiceLineDto[b].Id == invoiceUpdate.UpdateInvoiceLine[i].Id {
					notifier = 1
					break
				}
				notifier = 2
			}

			if notifier == 1 {
				resizeInvoiceLineDto = append(resizeInvoiceLineDto, invoiceUpdate.InvoiceLineDto[b])
				notifier = 0
				counter++
			} else if notifier == 2 {
				fmt.Println(invoiceUpdate.InvoiceLineDto[b].Id)
				err, count := invoiceLineService.Delete(invoiceUpdate.InvoiceLineDto[b].Id)
				if err != nil {
					return err, 0
				}
				if count == 0 {
					return errors.New("could not delete invoiceLine1"), 0
				}
				product := productService.FindById(invoiceUpdate.InvoiceLineDto[b].ItemId)
				product.QuantityOnHand = product.QuantityOnHand + invoiceUpdate.InvoiceLineDto[b].Quantity
				err, count = productService.Update(product)
				if err != nil {
					return err, 0
				}
				if count == 0 {
					return errors.New("could not update product2"), 0
				}
			}

		}

		// check for updates
		for i := 0; i < len(resizeInvoiceLineDto); i++ {
			if !reflect.DeepEqual(resizeInvoiceLineDto[i], invoiceUpdate.UpdateInvoiceLine[i]) {

				if invoiceUpdate.UpdateInvoiceLine[i].ItemId == invoiceUpdate.InvoiceLineDto[i].ItemId {

					product := productService.FindById(invoiceUpdate.UpdateInvoiceLine[i].Id)
					if invoiceUpdate.UpdateInvoiceLine[i].Quantity != invoiceUpdate.InvoiceLineDto[i].Quantity {
						if invoiceUpdate.InvoiceLineDto[i].Quantity > invoiceUpdate.UpdateInvoiceLine[i].Quantity {
							product.QuantityOnHand += invoiceUpdate.InvoiceLineDto[i].Quantity - invoiceUpdate.UpdateInvoiceLine[i].Quantity
							err, count := productService.Update(product)
							if err != nil {
								return err, 0
							}
							if count == 0 {
								return errors.New("could not update product3"), 0
							}
						} else {
							product.QuantityOnHand -= invoiceUpdate.UpdateInvoiceLine[i].Quantity - invoiceUpdate.InvoiceLineDto[i].Quantity
							err, count := productService.Update(product)
							if err != nil {
								return err, 0
							}
							if count == 0 {
								return errors.New("could not update product4"), 0
							}
						}
					}

				} else {
					product := productService.FindById(invoiceUpdate.InvoiceLineDto[i].ItemId)
					product.QuantityOnHand = invoiceUpdate.InvoiceLineDto[i].Quantity + product.QuantityOnHand
					err, count := productService.Update(product)
					if err != nil {
						return err, 0
					}
					if count == 0 {
						return errors.New("1could not update item"), 0
					}
					product1 := productService.FindById(invoiceUpdate.UpdateInvoiceLine[i].ItemId)
					quantityOnHand := product1.QuantityOnHand - invoiceUpdate.UpdateInvoiceLine[i].Quantity

					product1.QuantityOnHand = quantityOnHand
					err, count = productService.Update(product1)
					if err != nil {
						return err, 0
					}
					if count == 0 {
						return errors.New("2could not update item"), 0
					}
				}
				err, count := invoiceLineService.Update(invoiceUpdate.UpdateInvoiceLine[i])
				if err != nil {
					return err, 0
				}
				if count == 0 {
					return errors.New("could not update invoiceLine1"), 0
				}
			}
		}
		err, count := p.InvoiceRepository.Update(ToInvoice(InvoiceReceiver{Id: invoiceUpdate.Id, InvoiceUniqueId: invoiceUpdate.InvoiceUniqueId,
			InvoiceDate: invoiceUpdate.InvoiceDate, Customer: updateCustomerId, InvoiceTotal: totalPriceUpdated}))
		if err != nil {
			return err, 0
		}
		if count == 0 {
			return errors.New("could not update invoice1"), 0
		}
		if updateCustomerId == invoiceUpdate.OriginalCustomer {

			err, count = customerService.AddBalance(invoiceUpdate.OriginalCustomer, invoiceUpdate.InvoiceTotal)
			if err != nil {
				return err, 0
			}
			if count == 0 {
				return errors.New("could not return back money to customer"), 0
			}
			err, count = customerService.SubtractBalance(invoiceUpdate.OriginalCustomer, totalPriceUpdated)

			if err != nil {
				return err, 0
			}
			if count == 0 {
				return errors.New("could not subtract money from customer please notify admin"), 0
			}

		} else {
			err, count = customerService.SubtractBalance(invoiceUpdate.UpdateCustomer, totalPriceUpdated)

			if err != nil {
				return err, 0
			}
			if count == 0 {
				return errors.New("could not return money for previous customer please notify admin"), 0
			}
			err, count = customerService.AddBalance(invoiceUpdate.OriginalCustomer, invoiceUpdate.InvoiceTotal)

			if err != nil {
				return err, 0
			}
			if count == 0 {
				return errors.New("could not subtract money from new customer please notify admin"), 0
			}

		}
		return nil, 1

	}
	return nil, 4
}

// check update and create invoiceLine
func (p *InvoiceService) updateAndCreateInvoice(invoiceUpdate InvoiceUpdate) (error, int64) {
	tax, err := strconv.ParseFloat(os.Getenv("TAX_RATE"), 64)
	if err != nil {
		return err, 0
	}
	addTax := 1 + tax/100
	threshold, err := strconv.ParseFloat(os.Getenv("TAX_THRESHOLD"), 64)
	if err != nil {
		return err, 0
	}
	updateCustomerId = 0
	totalPriceUpdated = 0
	customerService := InitCustomerApiReceiver.CustomerService
	invoiceLineService := InitInvoiceLineServiceReceiver
	productService := InitProductApiReceiver.ProductService
	// check available quantities
	for index, _ := range invoiceUpdate.UpdateInvoiceLine {
		product := productService.FindById(invoiceUpdate.UpdateInvoiceLine[index].ItemId)
		if product.QuantityOnHand-invoiceUpdate.UpdateInvoiceLine[index].Quantity < 0 {
			return errors.New("we do not have enough quantity of item:" + invoiceUpdate.UpdateInvoiceLine[index].ItemName), 0
		}
	}

	//update total price
	for index, _ := range invoiceUpdate.UpdateInvoiceLine {
		item := productService.FindById(invoiceUpdate.UpdateInvoiceLine[index].ItemId)
		totalPriceUpdated += item.Price * invoiceUpdate.UpdateInvoiceLine[index].Quantity
	}
	if totalPriceUpdated > threshold {
		totalPriceUpdated = totalPriceUpdated * addTax
		totalPriceUpdated = math.Round(totalPriceUpdated*100) / 100
	}
	// check for customer update
	if invoiceUpdate.UpdateCustomer != invoiceUpdate.OriginalCustomer {
		customer, err, count := customerService.GetCustomerById(invoiceUpdate.UpdateCustomer)
		if err != nil {
			return err, 0
		}
		if count == 0 {
			return errors.New("new customer not found"), 0
		}
		updateCustomerId = customer.Id

		if customer.Balance-totalPriceUpdated < 0 {
			return errors.New("updated customer has not enough balance"), 0
		}
	} else {
		updateCustomerId = invoiceUpdate.OriginalCustomer

		customer, err, count := customerService.GetCustomerById(invoiceUpdate.OriginalCustomer)
		customerBalance := customer.Balance + invoiceUpdate.InvoiceTotal
		if err != nil {
			return err, 0
		}
		if count == 0 {
			return errors.New("could not find data for old customer"), 0
		}
		if customerBalance-totalPriceUpdated < 0 {
			return errors.New("customer did not have enough money"), 0
		}
	}

	if len(invoiceUpdate.UpdateInvoiceLine) > len(invoiceUpdate.InvoiceLineDto) {
		// check for create invoiceLines and update

		for b := 0; b < len(invoiceUpdate.UpdateInvoiceLine); b++ {
			// create new invoiceLine
			if invoiceUpdate.UpdateInvoiceLine[b].Id == 0 {
				err := invoiceLineService.Create(invoiceUpdate.UpdateInvoiceLine[b])
				if err != nil {
					return errors.New("can not create new invoiceLine2"), 0

				}
				product := productService.FindById(invoiceUpdate.UpdateInvoiceLine[b].ItemId)
				product.QuantityOnHand = product.QuantityOnHand - invoiceUpdate.UpdateInvoiceLine[b].Quantity
				err, count := productService.Update(product)
				if err != nil {
					return err, 0
				}
				if count == 0 {
					return errors.New("could not update product5"), 0
				}

			} else {
				for i := 0; i < len(invoiceUpdate.InvoiceLineDto); i++ {
					if invoiceUpdate.UpdateInvoiceLine[b].Id == invoiceUpdate.InvoiceLineDto[i].Id {
						if !reflect.DeepEqual(invoiceUpdate.InvoiceLineDto[i], invoiceUpdate.UpdateInvoiceLine) {
							if invoiceUpdate.UpdateInvoiceLine[b].ItemId == invoiceUpdate.InvoiceLineDto[i].ItemId {
								if invoiceUpdate.UpdateInvoiceLine[b].Quantity > invoiceUpdate.InvoiceLineDto[i].Quantity {
									product := productService.FindById(invoiceUpdate.InvoiceLineDto[b].ItemId)
									product.QuantityOnHand = product.QuantityOnHand - (invoiceUpdate.InvoiceLineDto[b].Quantity - invoiceUpdate.InvoiceLineDto[i].Quantity)
									err, count := productService.Update(product)
									if err != nil {
										return err, 0
									}
									if count == 0 {
										return errors.New("could not update product6"), 0
									}
									err, count = invoiceLineService.Update(invoiceUpdate.UpdateInvoiceLine[b])
									if err != nil {
										return err, 0
									}
									if count == 0 {
										return errors.New("could not update invoiceLine2"), 0
									}
								}
								if invoiceUpdate.UpdateInvoiceLine[b].Quantity < invoiceUpdate.InvoiceLineDto[i].Quantity {
									product := productService.FindById(invoiceUpdate.InvoiceLineDto[b].ItemId)
									product.QuantityOnHand = product.QuantityOnHand + (invoiceUpdate.InvoiceLineDto[i].Quantity - invoiceUpdate.InvoiceLineDto[b].Quantity)
									err, count := productService.Update(product)
									if err != nil {
										return err, 0
									}
									if count == 0 {
										return errors.New("could not update product7"), 0
									}
									err, count = invoiceLineService.Update(invoiceUpdate.UpdateInvoiceLine[b])
									if err != nil {
										return err, 0
									}
									if count == 0 {
										return errors.New("could not update invoiceLine3"), 0
									}
								}

							} else {

								product := productService.FindById(invoiceUpdate.InvoiceLineDto[i].ItemId)
								product.QuantityOnHand = invoiceUpdate.InvoiceLineDto[i].Quantity + product.QuantityOnHand
								err, count := productService.Update(product)
								if err != nil {
									return err, 0
								}
								if count == 0 {
									return errors.New("3could not update item"), 0
								}
								product1 := productService.FindById(invoiceUpdate.UpdateInvoiceLine[b].ItemId)
								quantityOnHand := product1.QuantityOnHand - invoiceUpdate.UpdateInvoiceLine[b].Quantity

								product1.QuantityOnHand = quantityOnHand
								err, count = productService.Update(product1)
								if err != nil {
									return err, 0
								}
								if count == 0 {
									return errors.New("4could not update item"), 0
								}
								err, count = invoiceLineService.Update(invoiceUpdate.UpdateInvoiceLine[b])
								if err != nil {
									return err, 0
								}
								if count == 0 {
									return errors.New("could not update invoiceLine4"), 0
								}

							}

						}
						//break
					}

				}
			}

		}
		// check for delete element
		notifier := 1
		for a := 0; a < len(invoiceUpdate.InvoiceLineDto); a++ {
			for b := 0; b < len(invoiceUpdate.UpdateInvoiceLine); b++ {
				if invoiceUpdate.UpdateInvoiceLine[b].Id == invoiceUpdate.InvoiceLineDto[a].Id {
					notifier = 0
					break
				}
			}
			if notifier == 1 {
				notifier = 0
				err, count := invoiceLineService.Delete(invoiceUpdate.InvoiceLineDto[a].Id)
				if err != nil {
					return err, 0
				}
				if count == 0 {
					return errors.New("could not delete invoiceLine2"), 0
				}
				product := productService.FindById(invoiceUpdate.InvoiceLineDto[a].ItemId)
				product.QuantityOnHand = product.QuantityOnHand + invoiceUpdate.InvoiceLineDto[a].Quantity
				err, count = productService.Update(product)
				if err != nil {
					return err, 0
				}
				if count == 0 {
					return errors.New("could not update product8"), 0
				}
			}

		}

		err, count := p.InvoiceRepository.Update(ToInvoice(InvoiceReceiver{Id: invoiceUpdate.Id, InvoiceUniqueId: invoiceUpdate.InvoiceUniqueId,
			InvoiceDate: invoiceUpdate.InvoiceDate, Customer: updateCustomerId, InvoiceTotal: totalPriceUpdated}))
		if err != nil {
			return err, 0
		}
		if count == 0 {
			return errors.New("could not update invoice3"), 0
		}
		if updateCustomerId == invoiceUpdate.OriginalCustomer {
			fmt.Println(invoiceUpdate.InvoiceTotal)
			err, count = customerService.AddBalance(invoiceUpdate.OriginalCustomer, invoiceUpdate.InvoiceTotal)
			if err != nil {
				return err, 0
			}
			if count == 0 {
				return errors.New("could not update customer to add balance"), 0
			}
			fmt.Println(totalPriceUpdated)
			err, count = customerService.SubtractBalance(invoiceUpdate.OriginalCustomer, totalPriceUpdated)
			if err != nil {
				return err, 0
			}
			if count == 0 {
				return errors.New("could not update customer to subtract balance"), 0
			}
		} else {
			err, count = customerService.SubtractBalance(invoiceUpdate.UpdateCustomer, totalPriceUpdated)
			if err != nil {
				return err, 0
			}
			if count == 0 {
				return errors.New("could not update customer to subtract balance 2"), 0
			}
			err, count = customerService.AddBalance(invoiceUpdate.OriginalCustomer, invoiceUpdate.InvoiceTotal)
			if err != nil {
				return err, 0
			}
			if count == 0 {
				return errors.New("could not update customer to add balance 2"), 0
			}

		}
		return nil, 1

	}
	return nil, 5
}

// only update and delete where lengths are equal
func (p *InvoiceService) updateInvoiceEqualLength(invoiceUpdate InvoiceUpdate) (error, int64) {
	updateCustomerId = 0
	tax, err := strconv.ParseFloat(os.Getenv("TAX_RATE"), 64)
	if err != nil {
		return err, 0
	}
	addTax := 1 + tax/100
	threshold, err := strconv.ParseFloat(os.Getenv("TAX_THRESHOLD"), 64)
	if err != nil {
		return err, 0
	}
	totalPriceUpdated = 0
	customerService := InitCustomerApiReceiver.CustomerService
	invoiceLineService := InitInvoiceLineServiceReceiver
	productService := InitProductApiReceiver.ProductService
	// check available quantities
	for index, _ := range invoiceUpdate.UpdateInvoiceLine {
		product := productService.FindById(invoiceUpdate.UpdateInvoiceLine[index].ItemId)
		if product.QuantityOnHand-invoiceUpdate.UpdateInvoiceLine[index].Quantity < 0 {
			return errors.New("we do not have enough quantity of item:" + invoiceUpdate.UpdateInvoiceLine[index].ItemName), 0
		}
	}

	//update total price
	for index, _ := range invoiceUpdate.UpdateInvoiceLine {
		item := productService.FindById(invoiceUpdate.UpdateInvoiceLine[index].ItemId)
		totalPriceUpdated += item.Price * invoiceUpdate.UpdateInvoiceLine[index].Quantity

	}
	if totalPriceUpdated > threshold {
		totalPriceUpdated = totalPriceUpdated * addTax
		totalPriceUpdated = math.Round(totalPriceUpdated*100) / 100
	}
	// check for customer update
	if invoiceUpdate.UpdateCustomer != invoiceUpdate.OriginalCustomer {
		customer, err, count := customerService.GetCustomerById(invoiceUpdate.UpdateCustomer)
		if err != nil {
			return err, 0
		}
		if count == 0 {
			return errors.New("new customer not found"), 0
		}
		updateCustomerId = customer.Id

		if customer.Balance-totalPriceUpdated < 0 {
			return errors.New("updated customer has not enough balance"), 0
		}
	} else {
		updateCustomerId = invoiceUpdate.OriginalCustomer

		customer, err, count := customerService.GetCustomerById(invoiceUpdate.OriginalCustomer)
		customerBalance := invoiceUpdate.InvoiceTotal + customer.Balance
		if err != nil {
			return err, 0
		}
		if count == 0 {
			return errors.New("could not find data for old customer"), 0
		}
		if customerBalance-totalPriceUpdated < 0 {
			return errors.New("customer did not have enough money"), 0
		}

	}
	fmt.Println("inside equal")
	// check for delete invoiceLines and update
	if len(invoiceUpdate.UpdateInvoiceLine) == len(invoiceUpdate.InvoiceLineDto) {
		for index, _ := range invoiceUpdate.UpdateInvoiceLine {
			fmt.Println("inside loop")
			if invoiceUpdate.UpdateInvoiceLine[index].Id == 0 {
				err, count := invoiceLineService.Delete(invoiceUpdate.InvoiceLineDto[index].Id)
				if err != nil {
					return err, 0
				}
				if count == 0 {
					return errors.New("could not delete invoiceLine3"), 0
				}
				product := productService.FindById(invoiceUpdate.InvoiceLineDto[index].ItemId)
				product.QuantityOnHand = product.QuantityOnHand + invoiceUpdate.InvoiceLineDto[index].Quantity
				err, count = productService.Update(product)
				if err != nil {
					return err, 0
				}
				if count == 0 {
					return errors.New("could not update product9"), 0
				}
				err = invoiceLineService.Create(invoiceUpdate.UpdateInvoiceLine[index])
				if err != nil {
					return errors.New("can not create new invoiceLine3"), 0

				}
				product1 := productService.FindById(invoiceUpdate.UpdateInvoiceLine[index].ItemId)
				product1.QuantityOnHand = product1.QuantityOnHand - invoiceUpdate.UpdateInvoiceLine[index].Quantity
				err, count = productService.Update(product)
				if err != nil {
					return err, 0
				}
				if count == 0 {
					return errors.New("could not update product10"), 0
				}
			} else if !reflect.DeepEqual(invoiceUpdate.UpdateInvoiceLine[index], invoiceUpdate.InvoiceLineDto[index]) {

				if invoiceUpdate.UpdateInvoiceLine[index].ItemId == invoiceUpdate.InvoiceLineDto[index].ItemId {

					if invoiceUpdate.UpdateInvoiceLine[index].Quantity > invoiceUpdate.InvoiceLineDto[index].Quantity {
						product := productService.FindById(invoiceUpdate.InvoiceLineDto[index].ItemId)
						product.QuantityOnHand = product.QuantityOnHand - (invoiceUpdate.UpdateInvoiceLine[index].Quantity - invoiceUpdate.InvoiceLineDto[index].Quantity)
						err, count := productService.Update(product)
						if err != nil {
							return err, 0
						}
						if count == 0 {
							return errors.New("could not update product11"), 0
						}
						err, count = invoiceLineService.Update(invoiceUpdate.UpdateInvoiceLine[index])
						if err != nil {
							return err, 0
						}
						if count == 0 {
							return errors.New("could not update invoiceLine6"), 0
						}
					}
					if invoiceUpdate.UpdateInvoiceLine[index].Quantity < invoiceUpdate.InvoiceLineDto[index].Quantity {
						product := productService.FindById(invoiceUpdate.InvoiceLineDto[index].ItemId)
						product.QuantityOnHand = product.QuantityOnHand + (invoiceUpdate.InvoiceLineDto[index].Quantity - invoiceUpdate.UpdateInvoiceLine[index].Quantity)
						err, count := productService.Update(product)
						if err != nil {
							return err, 0
						}
						if count == 0 {
							return errors.New("could not update product12"), 0
						}
						err, count = invoiceLineService.Update(invoiceUpdate.UpdateInvoiceLine[index])
						if err != nil {
							return err, 0
						}
						if count == 0 {
							return errors.New("could not update invoiceLine7"), 0
						}
					}

				} else {

					product := productService.FindById(invoiceUpdate.InvoiceLineDto[index].ItemId)
					product.QuantityOnHand = invoiceUpdate.InvoiceLineDto[index].Quantity + product.QuantityOnHand
					err, count := productService.Update(product)
					if err != nil {
						return err, 0
					}
					if count == 0 {
						return errors.New("5could not update item"), 0
					}
					product1 := productService.FindById(invoiceUpdate.UpdateInvoiceLine[index].ItemId)
					quantityOnHand := product1.QuantityOnHand - invoiceUpdate.UpdateInvoiceLine[index].Quantity

					product1.QuantityOnHand = quantityOnHand
					err, count = productService.Update(product1)
					if err != nil {
						return err, 0
					}
					if count == 0 {
						return errors.New("6could not update item"), 0
					}
					err, count = invoiceLineService.Update(invoiceUpdate.UpdateInvoiceLine[index])
					if err != nil {
						return err, 0
					}
					if count == 0 {
						return errors.New("could not update invoiceLine8"), 0
					}

				}
			}
		}
		err, count := p.InvoiceRepository.Update(ToInvoice(InvoiceReceiver{Id: invoiceUpdate.Id, InvoiceUniqueId: invoiceUpdate.InvoiceUniqueId,
			InvoiceDate: invoiceUpdate.InvoiceDate, Customer: updateCustomerId, InvoiceTotal: totalPriceUpdated}))
		if err != nil {
			return err, 0
		}
		if count == 0 {
			return errors.New("could not update invoice10"), 0
		}
		if updateCustomerId == invoiceUpdate.OriginalCustomer {

			err, count = customerService.AddBalance(invoiceUpdate.OriginalCustomer, invoiceUpdate.InvoiceTotal)
			if err != nil {
				return err, 0
			}
			if count == 0 {
				return errors.New("could not update customer to add money2"), 0
			}
			err, count = customerService.SubtractBalance(invoiceUpdate.OriginalCustomer, totalPriceUpdated)
			if err != nil {
				return err, 0
			}
			if count == 0 {
				return errors.New("could not update customer to subtract money2"), 0
			}

		} else {
			err, count = customerService.SubtractBalance(invoiceUpdate.UpdateCustomer, totalPriceUpdated)
			if err != nil {
				return err, 0
			}
			if count == 0 {
				return errors.New("could not update customer to subtract money 4"), 0
			}
			err, count = customerService.AddBalance(invoiceUpdate.OriginalCustomer, invoiceUpdate.InvoiceTotal)

			if err != nil {
				return err, 0
			}
			if count == 0 {
				return errors.New("could not update customer to Add money 4"), 0
			}

		}

	}
	return nil, 0
}
