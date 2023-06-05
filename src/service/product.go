package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"tokopedia-integration-service/config"
	"tokopedia-integration-service/src/model"
	"tokopedia-integration-service/src/util"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/segmentio/kafka-go"
)

type ProductService interface {
	ConsumeProductMessages()
}

type productService struct{}

func NewProductService() ProductService {
	return &productService{}
}

func (ps *productService) ConsumeProductMessages() {
	// Set up the Kafka reader for product topic
	cfg := config.Get()
	config := kafka.ReaderConfig{
		Brokers: []string{fmt.Sprintf("%s:%s", cfg.KafkaHost, cfg.KafkaPort)},
		Topic:   cfg.KafkaProductTopic,
		GroupID: cfg.KafkaProductConsumerGroup,
	}
	reader := kafka.NewReader(config)

	// Continuously read messages from Kafka
	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Println("Error reading message from Kafka:", err.Error())
			continue
		}

		// Change kafka message from byte to struct
		var productMessage *model.ProductMessage
		err = json.Unmarshal(msg.Value, &productMessage)
		if err != nil {
			fmt.Println("Can't unmarshal the kafka message")
			continue
		}

		if productMessage.Method == model.CREATE {
			createProductBody := util.ConvertProductToCreateProductRequest(productMessage)
			url := cfg.TokopediaURL + "fs/1/create?shop_id=1"
			resp, _ := SendPostRequest(createProductBody, url)

			isFailData, failDataRow := IsFailData(resp)
			if isFailData {
				fmt.Println("Failed to send HTTP CREATE Request with error: " + failDataRow)
			} else if IsFailStatus(resp) {
				fmt.Println("Failed to send HTTP CREATE Request with error: " + resp.Status)
			} else {
				// Get Tokopedia Product ID and create a request to omnichannel backend
				productResponse := util.ConvertResponseToProductResponse(resp.Body)
				productID := productResponse.Data.SuccessRowsData[0].ProductID
				updateProductIdRequest := util.ConvertProductIdToUpdateProductIdRequest(productID)
				url = cfg.OmnichannelURL + strconv.Itoa(productID)
				_, _ = SendPostRequest(updateProductIdRequest, url)

				fmt.Printf("Successfully created a new product with id: %d\n", productID)
			}

		} else if productMessage.Method == model.UPDATE {
			updateProductBody := util.ConvertProductToUpdateProductRequest(productMessage)
			url := cfg.TokopediaURL + "fs/1/edit?shop_id=1"
			resp, _ := SendPatchRequest(updateProductBody, url)

			isFailData, failDataRow := IsFailData(resp)
			if isFailData {
				fmt.Println("Failed to send HTTP UPDATE Request with error: " + failDataRow)
			} else if IsFailStatus(resp) {
				fmt.Println("Failed to send HTTP UPDATE Request with error: " + resp.Status)
			} else {
				fmt.Printf("Successfully updated product with id: %d\n", productMessage.TokopediaID)
			}

		} else { // productMessage.Method == model.DELETE
			deleteProductBody := util.ConvertProductToDeleteProductRequest(productMessage)
			url := cfg.TokopediaURL + "fs/1/delete?shop_id=1"
			resp, _ := SendPostRequest(deleteProductBody, url)

			isFailData, failDataRow := IsFailData(resp)
			if isFailData {
				fmt.Println("Failed to send HTTP DELETE Request with error: " + failDataRow)
			} else if IsFailStatus(resp) {
				fmt.Println("Failed to send HTTP DELETE Request with error: " + resp.Status)
			} else {
				fmt.Printf("Successfully deleted product with id: %d\n", productMessage.TokopediaID)
			}
		}
	}
}

func SendPostRequest(body *bytes.Buffer, url string) (*http.Response, error) {
	retryClient := retryablehttp.NewClient()
	resp, err := retryClient.Post(url, "application/json", body)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func SendPatchRequest(body *bytes.Buffer, url string) (*http.Response, error) {
	req, _ := retryablehttp.NewRequest("PATCH", url, body)
	req.Header.Set("Content-Type", "application/json")
	retryClient := retryablehttp.NewClient()
	resp, err := retryClient.Do(req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func IsFailData(resp *http.Response) (bool, string) {
	productResponse := util.ConvertResponseToProductResponse(resp.Body)
	failData := productResponse.Data.FailData

	return failData > 0, productResponse.Data.FailedRowsData[0].Error[0]
}

func IsFailStatus(resp *http.Response) bool {
	return (resp.StatusCode < 200 || resp.StatusCode > 299)
}
