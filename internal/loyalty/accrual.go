package loyalty

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Galish/loyalty-system/internal/config"
	"github.com/Galish/loyalty-system/internal/logger"
	repo "github.com/Galish/loyalty-system/internal/repository"
)

const (
	limiterInterval time.Duration = 1 * time.Second
	maxAttempts     uint          = 10
)

type accrualClient struct {
	cfg       *config.Config
	repo      repo.LoyaltyRepository
	requestCh chan *accrualRequest
}

type accrualRequest struct {
	order    OrderNumber
	user     string
	attempts uint
}

type Accrual struct {
	ID     OrderNumber `json:"order"`
	Status Status      `json:"status"`
	Value  float32     `json:"accrual"`
}

func (r *accrualRequest) isAttemptsExceeded() bool {
	return r.attempts >= maxAttempts-1
}

func newAccrualClient(repo repo.LoyaltyRepository, cfg *config.Config) *accrualClient {
	client := &accrualClient{
		cfg:       cfg,
		repo:      repo,
		requestCh: make(chan *accrualRequest),
	}

	go client.run()

	return client
}

func (c *accrualClient) run() {
	limiter := newLimiter(limiterInterval)

	for req := range c.requestCh {
		<-limiter.C

		go func(req *accrualRequest) {
			accrual, err := c.getOrderAccrual(req.order)
			if err != nil || !accrual.Status.isFinal() && !req.isAttemptsExceeded() {
				c.retry(req)
				return
			}

			c.applyOrderAccrual(accrual, req.user)
		}(req)
	}

	limiter.Close()
}

func (c *accrualClient) newOrder(order *Order) {
	c.requestCh <- &accrualRequest{
		order:    order.ID,
		user:     order.User,
		attempts: 0,
	}
}

func (c *accrualClient) retry(req *accrualRequest) {
	if !req.isAttemptsExceeded() {
		c.requestCh <- &accrualRequest{
			order:    req.order,
			attempts: req.attempts + 1,
		}
	}
}

func (c *accrualClient) getOrderAccrual(order OrderNumber) (*Accrual, error) {
	url := fmt.Sprintf("%s/api/orders/%s", c.cfg.AccrualAddr, order)

	logger.WithFields(logger.Fields{
		"URL": url,
	}).Info("API call to accrual service")

	resp, err := http.Get(url)
	if err != nil {
		logger.WithError(err).Debug("API call to accrual service failed")
		return nil, err
	}

	var accrual Accrual
	if err := json.NewDecoder(resp.Body).Decode(&accrual); err != nil {
		logger.WithError(err).Debug("cannot decode request JSON body")
		return nil, err
	}

	defer resp.Body.Close()

	logger.WithFields(logger.Fields{
		"accrual": accrual,
	}).Debug("accrual service response")

	return &accrual, nil
}

func (c *accrualClient) applyOrderAccrual(accrual *Accrual, user string) error {
	ctx := context.Background()

	err := c.repo.UpdateOrder(
		ctx,
		&repo.Order{
			ID:      accrual.ID.String(),
			Status:  string(accrual.Status),
			Accrual: accrual.Value,
		},
	)
	if err != nil {
		logger.WithError(err).Debug("unable to update order")
		return err
	}

	err = c.repo.Enroll(
		ctx,
		&repo.Enrollment{
			User:        user,
			Sum:         accrual.Value,
			ProcessedAt: time.Now(),
		},
	)
	if err != nil {
		logger.WithError(err).Debug("unable to update balance")
		return err
	}

	return nil
}

func (c *accrualClient) close() {
	close(c.requestCh)
}
