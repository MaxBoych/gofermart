package order_models

import "time"

type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
)

type OrderStorageData struct {
	OrderId   int64       `db:"order_id"`
	Number    string      `db:"number"`
	UserId    int64       `db:"user_id"`
	Status    OrderStatus `db:"status"`
	Accrual   *float64    `db:"accrual"`
	CreatedAt time.Time   `db:"created_at"`
	UpdatedAt time.Time   `db:"updated_at"`
}

type OrderResponseData struct {
	Number     string      `json:"number"`
	Status     OrderStatus `json:"status"`
	Accrual    *float64    `json:"accrual,omitempty"`
	UploadedAt string      `json:"uploaded_at"`
}

func OrderStorageToResponse(storageData []OrderStorageData) []OrderResponseData {
	response := make([]OrderResponseData, len(storageData))
	for i := 0; i < len(storageData); i++ {
		response[i] = OrderResponseData{
			Number:     storageData[i].Number,
			Status:     storageData[i].Status,
			Accrual:    storageData[i].Accrual,
			UploadedAt: storageData[i].CreatedAt.Format(time.RFC3339),
		}
	}
	return response
}
