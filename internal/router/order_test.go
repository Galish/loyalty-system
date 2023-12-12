package router

import (
	"bytes"
	"context"
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

func TestHandlerAddOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockLoyaltyRepository(ctrl)

	m.EXPECT().
		CreateOrder(
			gomock.Any(),
			gomock.Any(),
		).
		DoAndReturn(func(ctx context.Context, order *repo.Order) error {
			if order.ID == "277431151" {
				return repo.ErrOrderConflict
			}

			if order.ID == "2774311589" {
				return repo.ErrOrderExists
			}

			return nil
		}).
		AnyTimes()

	cfg := config.Config{SrvAddr: "8000"}
	loyaltyService := loyalty.NewService(m)

	authService := auth.NewService(nil, "yvdUuY)HSX}?&b")
	jwtToken, _ := authService.GenerateToken(&repo.User{ID: "123"})

	ts := httptest.NewServer(
		New(&cfg, authService, loyaltyService),
	)
	defer ts.Close()

	type request struct {
		path   string
		method string
		cookie *http.Cookie
		body   string
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
				"/api/user/order",
				"POST",
				&http.Cookie{
					Name:  auth.AuthCookieName,
					Value: jwtToken,
				},
				"12345678903",
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
				"/api/user/orders",
				"PUT",
				&http.Cookie{
					Name:  auth.AuthCookieName,
					Value: jwtToken,
				},
				"12345678903",
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
				"/api/user/orders",
				"POST",
				&http.Cookie{},
				"",
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
				"/api/user/orders",
				"POST",
				&http.Cookie{
					Name:  auth.AuthCookieName,
					Value: jwtToken,
				},
				"",
			},
			&want{
				http.StatusBadRequest,
				"invalid order number value\n",
				"text/plain; charset=utf-8",
			},
		},
		{
			"invalid order number",
			&request{
				"/api/user/orders",
				"POST",
				&http.Cookie{
					Name:  auth.AuthCookieName,
					Value: jwtToken,
				},
				"12345",
			},
			&want{
				http.StatusBadRequest,
				"invalid order number value\n",
				"text/plain; charset=utf-8",
			},
		},
		{
			"order added successfully",
			&request{
				"/api/user/orders",
				"POST",
				&http.Cookie{
					Name:  auth.AuthCookieName,
					Value: jwtToken,
				},
				"12345678903",
			},
			&want{
				http.StatusAccepted,
				"",
				"",
			},
		},
		{
			"order already exists",
			&request{
				"/api/user/orders",
				"POST",
				&http.Cookie{
					Name:  auth.AuthCookieName,
					Value: jwtToken,
				},
				"2774311589",
			},
			&want{
				http.StatusOK,
				"order has already been added",
				"text/plain; charset=utf-8",
			},
		},
		{
			"order added by another user",
			&request{
				"/api/user/orders",
				"POST",
				&http.Cookie{
					Name:  auth.AuthCookieName,
					Value: jwtToken,
				},
				"277431151",
			},
			&want{
				http.StatusConflict,
				"order has already been added by another user\n",
				"text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(
				tt.req.method,
				ts.URL+tt.req.path,
				bytes.NewBuffer([]byte(tt.req.body)),
			)
			require.NoError(t, err)

			req.AddCookie(tt.req.cookie)
			req.Header.Add("Content-Type", "text/plain")

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

func TestHandlerGetOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockLoyaltyRepository(ctrl)

	m.EXPECT().
		GetUserOrders(
			gomock.Any(),
			"12345",
		).
		Return([]*repo.Order{
			{
				ID:         "#11111",
				Status:     "NEW",
				Accrual:    100,
				UploadedAt: time.Date(2021, time.Month(2), 21, 1, 10, 30, 0, time.UTC),
				User:       "12345",
			},
		}, nil).
		AnyTimes()

	m.EXPECT().
		GetUserOrders(
			gomock.Any(),
			"12346",
		).
		Return([]*repo.Order{}, nil).
		AnyTimes()

	cfg := config.Config{SrvAddr: "8000"}
	loyaltyService := loyalty.NewService(m)

	authService := auth.NewService(nil, "yvdUuY)HSX}?&b")
	jwtToken, _ := authService.GenerateToken(&repo.User{ID: "12345"})
	jwtToken2, _ := authService.GenerateToken(&repo.User{ID: "12346"})

	ts := httptest.NewServer(
		New(&cfg, authService, loyaltyService),
	)
	defer ts.Close()

	type request struct {
		path   string
		method string
		cookie *http.Cookie
		body   string
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
				"/api/user/order",
				"GET",
				&http.Cookie{
					Name:  auth.AuthCookieName,
					Value: jwtToken,
				},
				"12345678903",
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
				"/api/user/orders",
				"PUT",
				&http.Cookie{
					Name:  auth.AuthCookieName,
					Value: jwtToken,
				},
				"12345678903",
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
				"/api/user/orders",
				"GET",
				&http.Cookie{},
				"",
			},
			&want{
				http.StatusUnauthorized,
				"",
				"",
			},
		},
		{
			"order list",
			&request{
				"/api/user/orders",
				"GET",
				&http.Cookie{
					Name:  auth.AuthCookieName,
					Value: jwtToken,
				},
				"",
			},
			&want{
				http.StatusOK,
				"[{\"number\":\"#11111\",\"status\":\"NEW\",\"accrual\":100,\"uploaded_at\":\"2021-02-21T01:10:30Z\",\"user_id\":\"12345\"}]\n",
				"application/json",
			},
		},
		{
			"empty list",
			&request{
				"/api/user/orders",
				"GET",
				&http.Cookie{
					Name:  auth.AuthCookieName,
					Value: jwtToken2,
				},
				"",
			},
			&want{
				http.StatusNoContent,
				"",
				"",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(
				tt.req.method,
				ts.URL+tt.req.path,
				bytes.NewBuffer([]byte(tt.req.body)),
			)
			require.NoError(t, err)

			req.AddCookie(tt.req.cookie)
			req.Header.Add("Content-Type", "text/plain")

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
