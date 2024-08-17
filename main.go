package main

import (
	"assignment/customers"
	"assignment/invoiceLines"
	"assignment/invoices"
	"assignment/middleware"
	"assignment/products"
	"assignment/suppliers"
	"assignment/users"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

func initDb() *gorm.DB {
	dsn := fmt.Sprintf("%v:%v@tcp(localhost:3306)/%v?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
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
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}
	db := initDb()
	route := gin.Default()
	// refresh token
	route.POST("/tokenRefresh", middleware.RefreshToken)
	route.POST("/logout", middleware.LogOut)
	// customer routes
	customerAPI := customers.InitCustomerApi(db)
	invoices.InitCustomerApiReceiver = customerAPI
	route.POST("customer/create", middleware.AuthMiddleWareMember(""), customerAPI.Create)
	route.PUT("customer/update", middleware.AuthMiddleWareMember(""), customerAPI.Update)
	route.PUT("customer/updateBalance", middleware.AuthMiddleWareMember(""), customerAPI.UpdateBalance)
	route.PUT("customer/addBalance", middleware.AuthMiddleWareMember(""), customerAPI.AddBalance)
	route.PUT("customer/subtractBalance", middleware.AuthMiddleWareMember(""), customerAPI.SubtractBalance)
	route.DELETE("customer/delete", middleware.AuthMiddleWareMember(""), customerAPI.Delete)
	route.GET("customers", middleware.AuthMiddleWareMember(""), customerAPI.FindAll)
	route.GET("customer/:id", middleware.AuthMiddleWareMember(""), customerAPI.FindById)

	// suppler routes
	supplierApi := suppliers.InitSupplierApi(db)
	route.POST("supplier/create", middleware.AuthMiddleWareMember(""), supplierApi.Create)
	route.PUT("supplier/update", middleware.AuthMiddleWareMember(""), supplierApi.Update)
	route.GET("suppliers", middleware.AuthMiddleWareMember(""), supplierApi.FindAll)
	route.GET("supplier/:name", middleware.AuthMiddleWareMember(""), supplierApi.FindByName)
	route.DELETE("supplier/delete", middleware.AuthMiddleWareMember(""), supplierApi.Delete)

	// product routes
	productApi := products.InitProductApi(db)
	invoices.InitProductApiReceiver = productApi
	route.POST("product/create", middleware.AuthMiddleWareMember(""), productApi.Create)
	route.PUT("product/update", middleware.AuthMiddleWareMember(""), productApi.Update)
	route.GET("products", middleware.AuthMiddleWareMember(""), productApi.FindAll)
	route.GET("product/:name", middleware.AuthMiddleWareMember(""), productApi.FindByName)
	route.DELETE("product/delete", middleware.AuthMiddleWareMember(""), productApi.Delete)
	route.PUT("/product/image", middleware.AuthMiddleWareMember(""), productApi.UpdateImage)

	// invoice routes
	invoiceApi := invoices.InitInvoiceApi(db)
	invoiceLines1 := invoiceLines.InitInvoiceLineService(db)
	invoices.InitInvoiceLineServiceReceiver = invoiceLines1
	route.POST("invoice/create", middleware.AuthMiddleWareMember(""), invoiceApi.Create)
	route.PUT("invoice/update", middleware.AuthMiddleWareMember(""), invoiceApi.Update)
	route.GET("invoices", middleware.AuthMiddleWareMember(""), invoiceApi.FindAll)
	route.DELETE("invoice/delete", middleware.AuthMiddleWareMember(""), invoiceApi.Delete)

	//user Routes
	userRoutes := users.InitUserApi(db)
	route.POST("user/create", middleware.AuthMiddleWareMember("Admin"), userRoutes.Create)
	route.PUT("user/updateUsername", middleware.AuthMiddleWareMember(""), userRoutes.UpdateUserName)
	route.PUT("user/updatePassword", middleware.AuthMiddleWareMember(""), userRoutes.UpdatePassword)
	route.GET("users", middleware.AuthMiddleWareMember("Admin"), userRoutes.FindAll)
	route.GET("user/:userName", middleware.AuthMiddleWareMember("Admin"), userRoutes.FindByUserName)
	route.DELETE("user/delete", middleware.AuthMiddleWareMember("Admin"), userRoutes.Delete)
	route.POST("user/login", userRoutes.Login)
	err = route.Run("localhost:8080")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
