package model

type Method int

const (
	CREATE Method = iota
	UPDATE
	DELETE
)

type KafkaProductMessage struct {
	Method               Method  `json:"method"`
	ID                   string  `json:"id"`
	Name                 string  `json:"name"`
	Price                int     `json:"price"`
	Weight               float32 `json:"weight"`
	Stock                int     `json:"stock"`
	Image                string  `json:"image"`
	Description          string  `json:"description"`
	TokopediaProductID   int     `json:"tokopedia_product_id"`
	TokopediaFsID        int     `json:"tokopedia_fs_id"`
	TokopediaShopID      int     `json:"tokopedia_shop_id"`
	TokopediaBearerToken string  `json:"tokopedia_bearer_token"`
}

type KafkaErrorMessage struct {
	Method      string `json:"method"`
	Url         string `json:"url"`
	RequestBody string `json:"request_body"`
	Error       string `json:"error"`
	Status      string `json:"status"`
	RequestTime string `json:"request_time"`
}
