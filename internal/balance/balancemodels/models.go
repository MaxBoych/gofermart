package balancemodels

import "time"

type BalanceStorageData struct {
	BalanceID int64     `db:"balance_id"`
	UserID    int64     `db:"user_id"`
	Current   float64   `db:"current"`
	Withdrawn float64   `db:"withdrawn"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type BalanceResponseData struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

func BalanceStorageToResponse(balanceStorage BalanceStorageData) BalanceResponseData {
	return BalanceResponseData{
		Current:   balanceStorage.Current,
		Withdrawn: balanceStorage.Withdrawn,
	}
}

type WithdrawRequestData struct {
	Order  string  `json:"order" validate:"required"`
	Sum    float64 `json:"sum" validate:"required"`
	UserID int64   `json:"-"`
}

type BalanceChangeData struct {
	Action string
	Sum    float64
	UserID int64
}

func (c *BalanceChangeData) IsWithdraw() bool {
	return c.Action == "-"
}

type WithdrawStorageData struct {
	WithdrawID int64     `db:"withdraw_id"`
	Order      string    `db:"order"`
	Sum        float64   `db:"sum"`
	UserID     int64     `db:"user_id"`
	CreatedAt  time.Time `db:"created_at"`
}

type WithdrawResponseData struct {
	Order       string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}

func WithdrawStorageToResponse(storageData []WithdrawStorageData) []WithdrawResponseData {
	response := make([]WithdrawResponseData, len(storageData))
	for i := 0; i < len(storageData); i++ {
		response[i] = WithdrawResponseData{
			Order:       storageData[i].Order,
			Sum:         storageData[i].Sum,
			ProcessedAt: storageData[i].CreatedAt.Format(time.RFC3339),
		}
	}
	return response
}
