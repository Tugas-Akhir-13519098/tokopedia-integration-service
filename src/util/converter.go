package util

import (
	"bytes"
	"encoding/json"
	"io"
	"tokopedia-integration-service/src/model"
)

func ConvertProductToCreateProductRequest(pm *model.ProductMessage) *bytes.Buffer {
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
	body, _ := json.Marshal(product)
	responseBody := bytes.NewBuffer(body)

	return responseBody
}

func ConvertProductToUpdateProductRequest(pm *model.ProductMessage) *bytes.Buffer {
	product := model.UpdateProductRequest{
		ID:          pm.TokopediaID,
		Name:        pm.Name,
		Price:       pm.Price,
		Weight:      pm.Weight,
		Stock:       pm.Stock,
		Pictures:    []model.Picture{{FilePath: pm.Image}},
		Description: pm.Description,
	}
	body, _ := json.Marshal(product)
	responseBody := bytes.NewBuffer(body)

	return responseBody
}

func ConvertProductToDeleteProductRequest(pm *model.ProductMessage) *bytes.Buffer {
	product := model.DeleteProductRequest{
		ProductID: []int{pm.TokopediaID},
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
