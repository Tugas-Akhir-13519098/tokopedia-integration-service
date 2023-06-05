package model

type Method int

const (
	CREATE Method = iota
	UPDATE
	DELETE
)

type ProductMessage struct {
	Method      Method
	ID          string
	Name        string
	Price       int
	Weight      float32
	Stock       int
	Image       string
	Description string
	TokopediaID int
	ShopeeID    int
}

type Picture struct {
	FilePath string `json:"file_path" binding:"required"`
}

type CreateProductRequest struct {
	Name          string    `json:"name" binding:"required"`
	CategoryId    int       `json:"category_id" binding:"required"`
	PriceCurrency string    `json:"price_currency" binding:"required"`
	Price         int       `json:"price" binding:"required"`
	Status        string    `json:"status" binding:"required"`
	MinOrder      int       `json:"min_order" binding:"required"`
	Weight        float32   `json:"weight" binding:"required"`
	WeightUnit    string    `json:"weight_unit" binding:"required"`
	Condition     string    `json:"condition" binding:"required"`
	Stock         int       `json:"stock" binding:"required"`
	Pictures      []Picture `json:"pictures" binding:"required"`
	Description   string    `json:"description" binding:"required"`
}

type UpdateProductRequest struct {
	ID          int       `json:"id" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	Price       int       `json:"price" binding:"required"`
	Weight      float32   `json:"weight" binding:"required"`
	Stock       int       `json:"stock" binding:"required"`
	Pictures    []Picture `json:"pictures" binding:"required"`
	Description string    `json:"description" binding:"required"`
}

type DeleteProductRequest struct {
	ProductID []int `json:"product_id" binding:"required"`
}

type ProductResponse struct {
	Header ProductResponseHeader `json:"header" binding:"required"`
	Data   ProductResponseData   `json:"data" binding:"required"`
}

type ProductResponseHeader struct {
	ProcessTime float32 `json:"process_time" binding:"required"`
	Messages    string  `json:"messages" binding:"required"`
}

type ProductResponseData struct {
	TotalData       int                   `json:"total_data" binding:"required"`
	SucceessData    int                   `json:"succeed_rows" binding:"required"`
	FailData        int                   `json:"failed_rows" binding:"required"`
	SuccessRowsData []SuccessResponseData `json:"success_rows_data"`
	FailedRowsData  []FailedResponseData  `json:"failed_rows_data"`
}

type SuccessResponseData struct {
	ProductID int `json:"product_id"`
}

type FailedResponseData struct {
	ProductID   int      `json:"product_id"`
	ProductName string   `json:"product_name"`
	Error       []string `json:"error"`
}

type UpdateProductIdRequest struct {
	TokopediaProductID int `json:"tokopedia_product_id" binding:"required"`
}
