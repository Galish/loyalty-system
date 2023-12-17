package router

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Galish/loyalty-system/internal/auth"
	"github.com/Galish/loyalty-system/internal/config"
	"github.com/Galish/loyalty-system/internal/loyalty"
	repo "github.com/Galish/loyalty-system/internal/repository"
	"github.com/Galish/loyalty-system/internal/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlerGetBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockLoyaltyRepository(ctrl)

	m.EXPECT().
		GetUserBalance(
			gomock.Any(),
			"395fd5f4-964d-4135-9a55-fbf91c4a163b",
		).
		Return(&repo.Balance{
			Current:   155.01,
			Withdrawn: 25.88,
		}, nil).
		AnyTimes()

	m.EXPECT().
		GetUserBalance(
			gomock.Any(),
			"395fd5f4-964d-4135-9a55-fbf91c4a1614",
		).
		Return(nil, errors.New("user not found")).
		AnyTimes()

	cfg := config.Config{SrvAddr: "8000"}
	loyaltyService := loyalty.NewService(m, &cfg)

	authService := auth.NewService(nil, "yvdUuY)HSX}?&b")
	jwtToken, _ := authService.GenerateToken(&repo.User{ID: "395fd5f4-964d-4135-9a55-fbf91c4a163b"})
	jwtToken2, _ := authService.GenerateToken(&repo.User{ID: "395fd5f4-964d-4135-9a55-fbf91c4a1614"})

	ts := httptest.NewServer(
		New(&cfg, authService, loyaltyService),
	)
	defer ts.Close()

	type request struct {
		path   string
		method string
		cookie *http.Cookie
	}

	type want struct {
		statusCode  int
		body        string
		contentType string
	}

	tests := []struct {
		name string
		req  *request
		want *want
	}{
		{
			"invalid API endpoint",
			&request{
				"/api/user/balances",
				"GET",
				&http.Cookie{
					Name:  auth.AuthCookieName,
					Value: jwtToken,
				},
			},
			&want{
				http.StatusNotFound,
				"404 page not found\n",
				"text/plain; charset=utf-8",
			},
		},
		{
			"invalid request method",
			&request{
				"/api/user/balance",
				"PUT",
				&http.Cookie{
					Name:  auth.AuthCookieName,
					Value: jwtToken,
				},
			},
			&want{
				http.StatusMethodNotAllowed,
				"",
				"",
			},
		},
		{
			"unauthorized",
			&request{
				"/api/user/balance",
				"GET",
				&http.Cookie{},
			},
			&want{
				http.StatusUnauthorized,
				"",
				"",
			},
		},
		{
			"reading from repo error",
			&request{
				"/api/user/balance",
				"GET",
				&http.Cookie{
					Name:  auth.AuthCookieName,
					Value: jwtToken2,
				},
			},
			&want{
				http.StatusInternalServerError,
				"user not found\n",
				"text/plain; charset=utf-8",
			},
		},
		{
			"successful response",
			&request{
				"/api/user/balance",
				"GET",
				&http.Cookie{
					Name:  auth.AuthCookieName,
					Value: jwtToken,
				},
			},
			&want{
				http.StatusOK,
				"{\"current\":155.01,\"withdrawn\":25.88}\n",
				"application/json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(
				tt.req.method,
				ts.URL+tt.req.path,
				nil,
			)
			require.NoError(t, err)

			req.AddCookie(tt.req.cookie)

			client := &http.Client{}
			resp, err := client.Do(req)
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, resp.StatusCode)

			raw, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, resp.Header.Get("Content-Type"), tt.want.contentType)
			assert.Equal(t, tt.want.body, string(raw))

			err = resp.Body.Close()
			require.NoError(t, err)
		})
	}
}

func TestHandlerWithdraw(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockLoyaltyRepository(ctrl)

	m.EXPECT().
		Withdraw(
			gomock.Any(),
			gomock.Any(),
		).
		DoAndReturn(func(ctx context.Context, withdraw *repo.Withdrawal) error {
			if withdraw.Sum >= 700 {
				return errors.New("internal server error")
			}

			if withdraw.Sum >= 100 {
				return repo.ErrInsufficientFunds
			}

			return nil
		}).
		AnyTimes()

	cfg := config.Config{SrvAddr: "8000"}
	loyaltyService := loyalty.NewService(m, &cfg)

	authService := auth.NewService(nil, "yvdUuY)HSX}?&b")
	jwtToken, _ := authService.GenerateToken(&repo.User{ID: "395fd5f4-964d-4135-9a55-fbf91c4a163b"})

	ts := httptest.NewServer(
		New(&cfg, authService, loyaltyService),
	)
	defer ts.Close()

	type withdraw struct {
		Order string  `json:"order"`
		Sum   float32 `json:"sum"`
	}

	type request struct {
		path   string
		method string
		cookie *http.Cookie
		body   *withdraw
	}

	type want struct {
		statusCode  int
		body        string
		contentType string
	}

	tests := []struct {
		name string
		req  *request
		want *want
	}{
		{
			"invalid API endpoint",
			&request{
				"/api/user/balance/withdrawal",
				"POST",
				&http.Cookie{
					Name:  auth.AuthCookieName,
					Value: jwtToken,
				},
				&withdraw{
					Order: "2377225624",
					Sum:   751,
				},
			},
			&want{
				http.StatusNotFound,
				"404 page not found\n",
				"text/plain; charset=utf-8",
			},
		},
		{
			"invalid request method",
			&request{
				"/api/user/balance/withdraw",
				"PUT",
				&http.Cookie{
					Name:  auth.AuthCookieName,
					Value: jwtToken,
				},
				&withdraw{
					Order: "2377225624",
					Sum:   751,
				},
			},
			&want{
				http.StatusMethodNotAllowed,
				"",
				"",
			},
		},
		{
			"unauthorized",
			&request{
				"/api/user/balance/withdraw",
				"POST",
				&http.Cookie{},
				&withdraw{
					Order: "2377225624",
					Sum:   751,
				},
			},
			&want{
				http.StatusUnauthorized,
				"",
				"",
			},
		},
		{
			"empty request body",
			&request{
				"/api/user/balance/withdraw",
				"POST",
				&http.Cookie{
					Name:  auth.AuthCookieName,
					Value: jwtToken,
				},
				nil,
			},
			&want{
				http.StatusBadRequest,
				"cannot decode request JSON body\n",
				"text/plain; charset=utf-8",
			},
		},
		{
			"invalid order number",
			&request{
				"/api/user/balance/withdraw",
				"POST",
				&http.Cookie{
					Name:  auth.AuthCookieName,
					Value: jwtToken,
				},
				&withdraw{
					Order: "12345",
					Sum:   751,
				},
			},
			&want{
				http.StatusUnprocessableEntity,
				"invalid order number value\n",
				"text/plain; charset=utf-8",
			},
		},
		{
			"insufficient funds",
			&request{
				"/api/user/balance/withdraw",
				"POST",
				&http.Cookie{
					Name:  auth.AuthCookieName,
					Value: jwtToken,
				},
				&withdraw{
					Order: "2377225624",
					Sum:   151,
				},
			},
			&want{
				http.StatusPaymentRequired,
				"insufficient funds in the account\n",
				"text/plain; charset=utf-8",
			},
		},
		{
			"repo operation error",
			&request{
				"/api/user/balance/withdraw",
				"POST",
				&http.Cookie{
					Name:  auth.AuthCookieName,
					Value: jwtToken,
				},
				&withdraw{
					Order: "2377225624",
					Sum:   751,
				},
			},
			&want{
				http.StatusInternalServerError,
				"unable to withdraw funds\n",
				"text/plain; charset=utf-8",
			},
		},
		{
			"successful withdrawal of funds",
			&request{
				"/api/user/balance/withdraw",
				"POST",
				&http.Cookie{
					Name:  auth.AuthCookieName,
					Value: jwtToken,
				},
				&withdraw{
					Order: "2377225624",
					Sum:   51,
				},
			},
			&want{
				http.StatusOK,
				"",
				"",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reqBody []byte
			var err error

			if tt.req.body != nil {
				reqBody, err = json.Marshal(tt.req.body)
				require.NoError(t, err)
			}

			req, err := http.NewRequest(
				tt.req.method,
				ts.URL+tt.req.path,
				bytes.NewBuffer(reqBody),
			)
			require.NoError(t, err)

			req.Header.Add("Content-Type", "application/json")
			req.AddCookie(tt.req.cookie)

			client := &http.Client{}
			resp, err := client.Do(req)
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, resp.StatusCode)

			raw, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, resp.Header.Get("Content-Type"), tt.want.contentType)
			assert.Equal(t, tt.want.body, string(raw))

			err = resp.Body.Close()
			require.NoError(t, err)
		})
	}
}

func TestHandlerGetWithdrawals(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	timezone1 := time.FixedZone("UTC-8", -8*60*60)
	timezone2 := time.FixedZone("UTC+5", 5*60*60)

	m := mocks.NewMockLoyaltyRepository(ctrl)

	m.EXPECT().
		GetWithdrawals(
			gomock.Any(),
			"395fd5f4-964d-4135-9a55-fbf91c4a163b",
		).
		Return([]*repo.Withdrawal{
			{
				Order:       "277431151",
				User:        "395fd5f4-964d-4135-9a55-fbf91c4a163b",
				Sum:         500,
				ProcessedAt: time.Date(2023, time.Month(2), 21, 1, 10, 30, 0, time.UTC),
			},
			{
				Order:       "277431113",
				User:        "395fd5f4-964d-4135-9a55-fbf91c4a163b",
				Sum:         150,
				ProcessedAt: time.Date(2023, time.Month(5), 21, 1, 10, 30, 0, timezone1),
			},
			{
				Order:       "277431122",
				User:        "395fd5f4-964d-4135-9a55-fbf91c4a163b",
				Sum:         755,
				ProcessedAt: time.Date(2023, time.Month(6), 21, 1, 10, 30, 0, timezone2),
			},
		}, nil).
		AnyTimes()

	m.EXPECT().
		GetWithdrawals(
			gomock.Any(),
			"395fd5f4-964d-4135-9a55-fbf91c4a1614",
		).
		Return([]*repo.Withdrawal{}, nil).
		AnyTimes()

	m.EXPECT().
		GetWithdrawals(
			gomock.Any(),
			"395fd5f4-964d-4135-9a55-fbf91c4a1615",
		).
		Return([]*repo.Withdrawal{}, errors.New("error occurred")).
		AnyTimes()

	cfg := config.Config{SrvAddr: "8000"}
	loyaltyService := loyalty.NewService(m, &cfg)

	authService := auth.NewService(nil, "yvdUuY)HSX}?&b")
	jwtToken, _ := authService.GenerateToken(&repo.User{ID: "395fd5f4-964d-4135-9a55-fbf91c4a163b"})
	jwtToken2, _ := authService.GenerateToken(&repo.User{ID: "395fd5f4-964d-4135-9a55-fbf91c4a1614"})
	jwtToken3, _ := authService.GenerateToken(&repo.User{ID: "395fd5f4-964d-4135-9a55-fbf91c4a1615"})

	ts := httptest.NewServer(
		New(&cfg, authService, loyaltyService),
	)
	defer ts.Close()

	type request struct {
		path   string
		method string
		cookie *http.Cookie
	}

	type want struct {
		statusCode  int
		body        string
		contentType string
	}

	tests := []struct {
		name string
		req  *request
		want *want
	}{
		{
			"invalid API endpoint",
			&request{
				"/api/user/withdrawal",
				"GET",
				&http.Cookie{
					Name:  auth.AuthCookieName,
					Value: jwtToken,
				},
			},
			&want{
				http.StatusNotFound,
				"404 page not found\n",
				"text/plain; charset=utf-8",
			},
		},
		{
			"invalid request method",
			&request{
				"/api/user/withdrawals",
				"PUT",
				&http.Cookie{
					Name:  auth.AuthCookieName,
					Value: jwtToken,
				},
			},
			&want{
				http.StatusMethodNotAllowed,
				"",
				"",
			},
		},
		{
			"unauthorized",
			&request{
				"/api/user/withdrawals",
				"GET",
				&http.Cookie{},
			},
			&want{
				http.StatusUnauthorized,
				"",
				"",
			},
		},
		{
			"no content",
			&request{
				"/api/user/withdrawals",
				"GET",
				&http.Cookie{
					Name:  auth.AuthCookieName,
					Value: jwtToken2,
				},
			},
			&want{
				http.StatusNoContent,
				"",
				"",
			},
		},
		{
			"reading from repo error",
			&request{
				"/api/user/withdrawals",
				"GET",
				&http.Cookie{
					Name:  auth.AuthCookieName,
					Value: jwtToken3,
				},
			},
			&want{
				http.StatusInternalServerError,
				"unable to get withdrawals\n",
				"text/plain; charset=utf-8",
			},
		},
		{
			"successful response",
			&request{
				"/api/user/withdrawals",
				"GET",
				&http.Cookie{
					Name:  auth.AuthCookieName,
					Value: jwtToken,
				},
			},
			&want{
				http.StatusOK,
				"[{\"order\":\"277431151\",\"sum\":500,\"processed_at\":\"2023-02-21T01:10:30+00:00\"},{\"order\":\"277431113\",\"sum\":150,\"processed_at\":\"2023-05-21T01:10:30-08:00\"},{\"order\":\"277431122\",\"sum\":755,\"processed_at\":\"2023-06-21T01:10:30+05:00\"}]\n",
				"application/json",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(
				tt.req.method,
				ts.URL+tt.req.path,
				nil,
			)
			require.NoError(t, err)

			req.AddCookie(tt.req.cookie)

			client := &http.Client{}
			resp, err := client.Do(req)
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, resp.StatusCode)

			raw, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, resp.Header.Get("Content-Type"), tt.want.contentType)
			assert.Equal(t, tt.want.body, string(raw))

			err = resp.Body.Close()
			require.NoError(t, err)
		})
	}
}
