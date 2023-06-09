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
		var productMessage *model.ProductMessage
		err = json.Unmarshal(msg.Value, &productMessage)
		if err != nil {
			fmt.Println("Can't unmarshal the kafka message")
			continue
		}

		if productMessage.Method == model.CREATE {
			createProductBody := util.ConvertProductToCreateProductRequest(productMessage)
			url := cfg.TokopediaURL + "fs/1/create?shop_id=1"
			resp, _ := util.SendPostRequest(createProductBody, url)
			util.AfterHTTPRequestHandler(resp, "CREATE", productMessage.ID, cfg.OmnichannelURL)

		} else if productMessage.Method == model.UPDATE {
			updateProductBody := util.ConvertProductToUpdateProductRequest(productMessage)
			url := cfg.TokopediaURL + "fs/1/edit?shop_id=1"
			resp, _ := util.SendPatchRequest(updateProductBody, url)
			util.AfterHTTPRequestHandler(resp, "UPDATE", string(msg.Key), cfg.OmnichannelURL)

		} else { // productMessage.Method == model.DELETE
			deleteProductBody := util.ConvertProductToDeleteProductRequest(productMessage)
			url := cfg.TokopediaURL + "fs/1/delete?shop_id=1"
			resp, _ := util.SendPostRequest(deleteProductBody, url)
			util.AfterHTTPRequestHandler(resp, "DELETE", string(msg.Key), cfg.OmnichannelURL)
		}
	}
}
