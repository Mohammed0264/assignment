package main

import (
	"assignment/customers"
	"assignment/invoiceLines"
	"assignment/invoices"
	"assignment/products"
	"assignment/suppliers"
	"assignment/users"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func initDb() *gorm.DB {
	dsn := "root:root@tcp(localhost:3306)/assignment?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	er := db.AutoMigrate(&customers.Customer{}, &suppliers.Supplier{}, &products.Product{}, &invoices.Invoice{},
		&invoiceLines.InvoiceLine{}, &users.User{})

	if err != nil {
		panic("failed to connect database")
		return nil
	}
	if er != nil {
		fmt.Println(er.Error())
		return nil
	}
	fmt.Println("Successfully connected to database")
	return db
}

func main() {
	db := initDb()
	route := gin.Default()
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	// customer routes
	customerAPI := customers.InitCustomerApi(db)
	invoices.InitCustomerApiReceiver = customerAPI
	route.POST("customer/create", customerAPI.Create)
	route.PUT("customer/update", customerAPI.Update)
	route.PUT("customer/updateBalance", customerAPI.UpdateBalance)
	route.PUT("customer/addBalance", customerAPI.AddBalance)
	route.PUT("customer/subtractBalance", customerAPI.SubtractBalance)
	route.DELETE("customer/delete", customerAPI.Delete)
	route.GET("customers", customerAPI.FindAll)
	route.GET("customer/:id", customerAPI.FindById)

	// suppler routes
	supplierApi := suppliers.InitSupplierApi(db)
	route.POST("supplier/create", supplierApi.Create)
	route.PUT("supplier/update", supplierApi.Update)
	route.GET("suppliers", supplierApi.FindAll)
	route.GET("supplier/:name", supplierApi.FindByName)
	route.DELETE("supplier/delete", supplierApi.Delete)

	// product routes
	productApi := products.InitProductApi(db)
	invoices.InitProductApiReceiver = productApi
	route.POST("product/create", productApi.Create)
	route.PUT("product/update", productApi.Update)
	route.GET("products", productApi.FindAll)
	route.GET("product/:name", productApi.FindByName)
	route.DELETE("product/delete", productApi.Delete)
	route.PUT("/product/image", productApi.UpdateImage)

	// invoice routes
	invoiceApi := invoices.InitInvoiceApi(db)
	invoiceLines1 := invoiceLines.InitInvoiceLineService(db)
	invoices.InitInvoiceLineServiceReceiver = invoiceLines1
	route.POST("invoice/create", invoiceApi.Create)
	route.PUT("invoice/update", invoiceApi.Update)
	route.GET("invoices", invoiceApi.FindAll)
	route.DELETE("invoice/delete", invoiceApi.Delete)

	//user Routes
	userRoutes := users.InitUserApi(db)
	route.POST("user/create", userRoutes.Create)
	route.PUT("user/updateUsername", userRoutes.UpdateUserName)
	route.PUT("user/updatePassword", userRoutes.UpdatePassword)
	route.GET("users", userRoutes.FindAll)
	route.GET("user/:userName", userRoutes.FindByUserName)
	route.DELETE("user/delete", userRoutes.Delete)
	err = route.Run("localhost:8080")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
