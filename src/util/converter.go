package util

import (
	"bytes"
	"encoding/json"
	"io"
	"tokopedia-integration-service/src/model"
)

func ConvertProductToCreateProductRequest(pm *model.KafkaProductMessage) *bytes.Buffer {
	product := model.CreateProductRequest{
		Name:          pm.Name,
		CategoryId:    1,
		PriceCurrency: "IDR",
		Price:         pm.Price,
		Status:        "LIMITED",
		MinOrder:      1,
		Weight:        pm.Weight,
		WeightUnit:    "GR",
		Condition:     "NEW",
		Stock:         pm.Stock,
		Pictures:      []model.Picture{{FilePath: pm.Image}},
		Description:   pm.Description,
	}
	request := model.CreateRequest{
		Products: []model.CreateProductRequest{product},
	}
	body, _ := json.Marshal(request)
	responseBody := bytes.NewBuffer(body)

	return responseBody
}

func ConvertProductToUpdateProductRequest(pm *model.KafkaProductMessage) *bytes.Buffer {
	product := model.UpdateProductRequest{
		ID:          pm.TokopediaProductID,
		Name:        pm.Name,
		Price:       pm.Price,
		Weight:      pm.Weight,
		Stock:       pm.Stock,
		Pictures:    []model.Picture{{FilePath: pm.Image}},
		Description: pm.Description,
	}
	request := model.UpdateRequest{
		Products: []model.UpdateProductRequest{product},
	}
	body, _ := json.Marshal(request)
	responseBody := bytes.NewBuffer(body)

	return responseBody
}

func ConvertProductToDeleteProductRequest(pm *model.KafkaProductMessage) *bytes.Buffer {
	product := model.DeleteProductRequest{
		ProductID: []int{pm.TokopediaProductID},
	}
	body, _ := json.Marshal(product)
	responseBody := bytes.NewBuffer(body)

	return responseBody
}

func ConvertResponseToProductResponse(body io.ReadCloser) model.ProductResponse {
	respBody, _ := io.ReadAll(body)
	var productResponse model.ProductResponse
	_ = json.Unmarshal(respBody, &productResponse)

	return productResponse
}

func ConvertProductIdToUpdateProductIdRequest(productID int) *bytes.Buffer {
	request := model.UpdateProductIdRequest{
		TokopediaProductID: productID,
	}
	body, _ := json.Marshal(request)
	responseBody := bytes.NewBuffer(body)

	return responseBody
}

func ConvertToErrorMessage(method string, url string, req string, err string, status string, reqTime string) []byte {
	message := model.KafkaErrorMessage{
		Method:      method,
		Url:         url,
		RequestBody: req,
		Error:       err,
		Status:      status,
		RequestTime: reqTime,
	}
	messageByte, _ := json.Marshal(message)

	return messageByte
}
