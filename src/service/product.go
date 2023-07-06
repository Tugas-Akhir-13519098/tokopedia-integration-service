package service

import (
	"context"
	"encoding/json"
	"fmt"
	"tokopedia-integration-service/config"
	"tokopedia-integration-service/src/model"
	"tokopedia-integration-service/src/util"

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
		var productMessage *model.KafkaProductMessage
		err = json.Unmarshal(msg.Value, &productMessage)
		if err != nil {
			fmt.Println("Can't unmarshal the kafka message")
			continue
		}

		if productMessage.Method == model.CREATE {
			createProductBody := util.ConvertProductToCreateProductRequest(productMessage)
			url := cfg.TokopediaURL + fmt.Sprintf("fs/%d/create?shop_id=%d", productMessage.TokopediaFsID, productMessage.TokopediaShopID)
			resp, _ := util.SendPostRequest(createProductBody, url, productMessage.TokopediaBearerToken)
			util.AfterHTTPRequestHandler(createProductBody.String(), resp, "CREATE", "POST", productMessage.ID, url)

		} else if productMessage.Method == model.UPDATE {
			updateProductBody := util.ConvertProductToUpdateProductRequest(productMessage)
			url := cfg.TokopediaURL + fmt.Sprintf("fs/%d/edit?shop_id=%d", productMessage.TokopediaFsID, productMessage.TokopediaShopID)
			resp, _ := util.SendPatchRequest(updateProductBody, url, productMessage.TokopediaBearerToken)
			util.AfterHTTPRequestHandler(updateProductBody.String(), resp, "UPDATE", "PATCH", string(msg.Key), url)

		} else { // productMessage.Method == model.DELETE
			deleteProductBody := util.ConvertProductToDeleteProductRequest(productMessage)
			url := cfg.TokopediaURL + fmt.Sprintf("fs/%d/delete?shop_id=%d", productMessage.TokopediaFsID, productMessage.TokopediaShopID)
			resp, _ := util.SendPostRequest(deleteProductBody, url, productMessage.TokopediaBearerToken)
			util.AfterHTTPRequestHandler(deleteProductBody.String(), resp, "DELETE", "POST", string(msg.Key), url)
		}
	}
}
