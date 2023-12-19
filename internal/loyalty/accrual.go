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
	maxAttempts     uint          = 5
)

type accrualClient struct {
	cfg       *config.Config
	repo      repo.LoyaltyRepository
	requestCh chan *accrualRequest
}

type accrualRequest struct {
	order    *Order
	attempts uint
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

			c.applyOrderAccrual(accrual)
		}(req)
	}

	limiter.Close()
}

func (c *accrualClient) newOrder(order *Order) {
	c.requestCh <- &accrualRequest{
		order:    order,
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

func (c *accrualClient) getOrderAccrual(order *Order) (*Order, error) {
	url := fmt.Sprintf("%s/api/orders/%s", c.cfg.AccrualAddr, order.ID)

	logger.WithFields(logger.Fields{
		"URL": url,
	}).Info("API call to accrual service")

	resp, err := http.Get(url)
	if err != nil {
		logger.WithError(err).Debug("API call to accrual service failed")
		return nil, err
	}

	var payload Order
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		logger.WithError(err).Debug("cannot decode request JSON body")
		return nil, err
	}

	defer resp.Body.Close()

	payload.User = order.User

	return &payload, nil
}

func (c *accrualClient) applyOrderAccrual(order *Order) error {
	ctx := context.Background()

	err := c.repo.UpdateOrder(
		ctx,
		&repo.Order{
			ID:      order.ID.String(),
			Status:  string(order.Status),
			Accrual: order.Accrual,
		},
	)
	if err != nil {
		logger.WithError(err).Debug("unable to update order")
		return err
	}

	err = c.repo.Enroll(
		ctx,
		&repo.Enrollment{
			User:        order.User,
			Sum:         order.Accrual,
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
