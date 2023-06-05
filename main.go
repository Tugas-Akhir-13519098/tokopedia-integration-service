package main

import (
	"log"
	"tokopedia-integration-service/src/service"
)

func main() {
	log.Printf("Application is running")

	productService := service.NewProductService()
	productService.ConsumeProductMessages()
}
