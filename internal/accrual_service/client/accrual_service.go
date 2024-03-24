package client

import (
	"encoding/json"
	"github.com/MaxBoych/gofermart/internal/accrual_service/accrual_service_models"
	"github.com/MaxBoych/gofermart/internal/order/order_models"
	"github.com/MaxBoych/gofermart/internal/store/consts"
	"github.com/MaxBoych/gofermart/pkg/errs"
	"github.com/MaxBoych/gofermart/pkg/logger"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
)

type AccrualServiceClient struct {
	AccrualSystemAddress string
	Queue                chan accrual_service_models.AccrualRequestWithResponse
}

func NewAccrualServiceClient(
	accrualSystemAddress string,
	queue chan accrual_service_models.AccrualRequestWithResponse,
) *AccrualServiceClient {
	return &AccrualServiceClient{
		AccrualSystemAddress: accrualSystemAddress,
		Queue:                queue,
	}
}

func (c *AccrualServiceClient) Run() {
	go func() {
		for reqWithResp := range c.Queue {
			resp, err := http.DefaultClient.Do(reqWithResp.Request)
			if err != nil {
				logger.Log.Error("Error to do request to accrual service", zap.Error(err))
				continue
			}
			reqWithResp.Response <- resp
		}
	}()
}

func (c *AccrualServiceClient) SendRequest(order order_models.OrderStorageData) (*http.Response, error) {
	for {
		req, err := http.NewRequest(
			"GET",
			c.AccrualSystemAddress+"/api/orders/"+order.Number,
			nil,
		)
		if err != nil {
			logger.Log.Error("Error to create http.NewRequest", zap.Error(err))
			return nil, err
		}

		responseChan := make(chan *http.Response)
		reqWithResp := accrual_service_models.AccrualRequestWithResponse{
			Request:  req,
			Response: responseChan,
		}
		c.Queue <- reqWithResp

		resp := <-responseChan
		close(responseChan)
		logger.Log.Info("Status for order " + order.Number + " is " + resp.Status)

		switch resp.StatusCode {
		case http.StatusOK:
			return resp, nil
		case http.StatusTooManyRequests:
			time.Sleep(consts.SleepTime)
			continue
		case http.StatusNoContent:
			return nil, errs.HttpErrOrderNoContent
		default:
			return nil, errs.HttpErrInternal
		}
	}
}

func (c *AccrualServiceClient) HttpResponseToOrderAccrualResponse(resp *http.Response) (*accrual_service_models.AccrualOrderResponseData, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Log.Error("Error to read response body", zap.Error(err))
		return nil, err
	}

	var data accrual_service_models.AccrualOrderResponseData
	if err := json.Unmarshal(body, &data); err != nil {
		logger.Log.Error("Error to do unmarshal response body", zap.Error(err))
		return nil, err
	}

	return &data, nil
}
