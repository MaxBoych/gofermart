package client

import (
	"encoding/json"
	"github.com/MaxBoych/gofermart/internal/accrual_service/accrualservicemodels"
	"github.com/MaxBoych/gofermart/internal/order/ordermodels"
	"github.com/MaxBoych/gofermart/pkg/consts"
	"github.com/MaxBoych/gofermart/pkg/errs"
	"github.com/MaxBoych/gofermart/pkg/logger"
	"go.uber.org/zap"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"time"
)

type AccrualServiceClient struct {
	AccrualSystemAddress string
	Queue                chan accrualservicemodels.AccrualRequestWithResponse
}

func NewAccrualServiceClient(
	accrualSystemAddress string,
	queue chan accrualservicemodels.AccrualRequestWithResponse,
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
				reqWithResp.Response <- nil
				continue
			}
			logger.Log.Info("Status for order is " + resp.Status)
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				logger.Log.Error("Error to read response body", zap.Error(err))
				reqWithResp.Response <- nil
				continue
			}

			reqWithResp.Response <- &accrualservicemodels.ResponseData{
				Body:   body,
				Status: uint16(resp.StatusCode),
			}
			resp.Body.Close()
		}
	}()
}

func (c *AccrualServiceClient) SendRequest(order ordermodels.OrderStorageData) ([]byte, error) {
	for i := 0; i < 3; i++ {
		req, err := http.NewRequest(
			"GET",
			c.AccrualSystemAddress+"/api/orders/"+order.Number,
			nil,
		)
		if err != nil {
			logger.Log.Error("Error to create http.NewRequest", zap.Error(err))
			return nil, err
		}

		reqDump, _ := httputil.DumpRequestOut(req, true)
		logger.Log.Info("request data", zap.String("request", string(reqDump)))

		responseChan := make(chan *accrualservicemodels.ResponseData)
		reqWithResp := accrualservicemodels.AccrualRequestWithResponse{
			Request:  req,
			Response: responseChan,
		}
		c.Queue <- reqWithResp

		resp := <-responseChan
		close(responseChan)
		if resp == nil {
			continue
		}

		switch resp.Status {
		case http.StatusOK:
			return resp.Body, nil
		case http.StatusTooManyRequests:
			i--
			time.Sleep(consts.SleepTime)
			continue
		case http.StatusNoContent:
			return nil, errs.HTTPErrOrderNoContent
		default:
			return nil, errs.HTTPErrInternal
		}
	}

	return nil, errs.HTTPErrConnectionRefused
}

func (c *AccrualServiceClient) HTTPResponseToOrderAccrualResponse(body []byte) (*accrualservicemodels.AccrualOrderResponseData, error) {
	logger.Log.Info("Response body content", zap.String("content", string(body)))

	var data accrualservicemodels.AccrualOrderResponseData
	if err := json.Unmarshal(body, &data); err != nil {
		logger.Log.Error("Error to do unmarshal response body", zap.Error(err))
		return nil, err
	}

	return &data, nil
}

func (c *AccrualServiceClient) CheckConnection() bool {
	_, err := net.Dial("tcp", c.AccrualSystemAddress)
	if err != nil {
		logger.Log.Error("Error to connect to accrual service", zap.Error(err))
		return false
	}
	return true
}
