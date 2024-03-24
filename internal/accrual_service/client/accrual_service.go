package accrual_service

import "net/http"

type AccrualServiceClient struct {
	queue chan *http.Request
	errs  chan error
}

func NewAccrualServiceClient(queue chan *http.Request, errs chan error) *AccrualServiceClient {
	return &AccrualServiceClient{
		queue: queue,
		errs:  errs,
	}
}

func (c *AccrualServiceClient) DoRequest(order string) {

}
