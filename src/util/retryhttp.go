package util

import (
	"bytes"
	"fmt"
	"net/http"
	"tokopedia-integration-service/src/model"

	"github.com/hashicorp/go-retryablehttp"
)

func SendPostRequest(body *bytes.Buffer, url string) (*http.Response, error) {
	retryClient := NewRetryClient()
	resp, err := retryClient.Post(url, "application/json", body)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func SendPatchRequest(body *bytes.Buffer, url string) (*http.Response, error) {
	retryClient := NewRetryClient()
	req, _ := retryablehttp.NewRequest("PATCH", url, body)
	req.Header.Set("Content-Type", "application/json")
	resp, err := retryClient.Do(req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func SendPutRequest(body *bytes.Buffer, url string) (*http.Response, error) {
	retryClient := NewRetryClient()
	req, _ := retryablehttp.NewRequest("PUT", url, body)
	req.Header.Set("Content-Type", "application/json")
	resp, err := retryClient.Do(req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func CustomErrorHandler(resp *http.Response, err error, numTries int) (*http.Response, error) {
	resp.Body.Close()
	return resp, err
}

func NewRetryClient() *retryablehttp.Client {
	retryClient := retryablehttp.NewClient()
	retryClient.ErrorHandler = CustomErrorHandler

	return retryClient
}

func AfterHTTPRequestHandler(resp *http.Response, method string, productID string, omnichannelURL string) {
	productResponse := ConvertResponseToProductResponse(resp.Body)
	IsFailResponse, failDataRow := IsFailResponse(resp, productResponse)
	if IsFailResponse {
		fmt.Printf("Failed to send HTTP %s Request with status: %s and error: %s\n", method, resp.Status, failDataRow)
	} else {
		if method == "CREATE" {
			// Get Tokopedia Product ID and create a request to omnichannel backend
			tokopediaID := productResponse.Data.SuccessRowsData[0].ProductID
			updateProductIdRequest := ConvertProductIdToUpdateProductIdRequest(tokopediaID)
			url := omnichannelURL + productID
			_, _ = SendPutRequest(updateProductIdRequest, url)

			fmt.Printf("Successfully CREATE a product with id: %s in Tokopedia\n", productID)
		} else {
			fmt.Printf("Successfully %s product with id: %s\n", method, productID)
		}
	}
}

func IsFailResponse(resp *http.Response, productResponse model.ProductResponse) (bool, string) {
	failData := productResponse.Data.FailData
	isFailStatus := (resp.StatusCode < 200 || resp.StatusCode > 299)

	if failData > 0 {
		return true, productResponse.Data.FailedRowsData[0].Error[0]
	} else if isFailStatus {
		return true, ""
	} else {
		return false, ""
	}
}
