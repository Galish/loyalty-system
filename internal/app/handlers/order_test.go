package handlers

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	repo "github.com/Galish/loyalty-system/internal/app/adapters/repository"
	repoMocks "github.com/Galish/loyalty-system/internal/app/adapters/repository/mocks"
	"github.com/Galish/loyalty-system/internal/app/entity"
	usecaseMocks "github.com/Galish/loyalty-system/internal/app/usecase/mocks"
	"github.com/Galish/loyalty-system/internal/app/usecase/order"
	"github.com/Galish/loyalty-system/internal/app/usecase/user"
	"github.com/Galish/loyalty-system/internal/auth"
	"github.com/Galish/loyalty-system/internal/config"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlerAddOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderMock := repoMocks.NewMockOrderRepository(ctrl)

	orderMock.EXPECT().
		CreateOrder(
			gomock.Any(),
			gomock.Any(),
		).
		DoAndReturn(func(ctx context.Context, order *entity.Order) error {
			if order.ID == "277431151" {
				return repo.ErrOrderConflict
			}

			if order.ID == "2774311589" {
				return repo.ErrOrderExists
			}

			return nil
		}).
		AnyTimes()

	accrualMock := usecaseMocks.NewMockAccrualUseCase(ctrl)

	accrualMock.EXPECT().GetAccrual(
		gomock.Any(),
		gomock.Eq(&entity.Order{
			ID:   "12345678903",
			User: "395fd5f4-964d-4135-9a55-fbf91c4a163b",
		}),
	)

	cfg := config.Config{
		SrvAddr:   "8000",
		SecretKey: "yvdUuY)HSX}?&b",
	}

	orderUseCase := order.New(orderMock)
	userUseCase := user.New(nil, cfg.SecretKey)

	jwtToken, _ := auth.GenerateToken(
		cfg.SecretKey,
		&entity.User{ID: "395fd5f4-964d-4135-9a55-fbf91c4a163b"},
	)

	ts := httptest.NewServer(
		NewRouter(
			&cfg,
			NewHandler(&cfg, accrualMock, nil, orderUseCase, userUseCase),
		),
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
				http.StatusUnprocessableEntity,
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
				http.StatusUnprocessableEntity,
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

	m := repoMocks.NewMockOrderRepository(ctrl)

	timezone1 := time.FixedZone("UTC-8", -8*60*60)
	timezone2 := time.FixedZone("UTC+5", 5*60*60)

	m.EXPECT().
		UserOrders(
			gomock.Any(),
			"395fd5f4-964d-4135-9a55-fbf91c4a163b",
		).
		Return([]*entity.Order{
			{
				ID:         "2774311589",
				Status:     "NEW",
				Accrual:    0,
				UploadedAt: time.Date(2023, time.Month(3), 11, 1, 10, 30, 0, timezone1),
				User:       "395fd5f4-964d-4135-9a55-fbf91c4a163b",
			},
			{
				ID:         "12345678903",
				Status:     "PROCESSED",
				Accrual:    100,
				UploadedAt: time.Date(2023, time.Month(2), 21, 1, 10, 30, 0, time.UTC),
				User:       "395fd5f4-964d-4135-9a55-fbf91c4a163b",
			},
			{
				ID:         "252576137",
				Status:     "INVALID",
				Accrual:    0,
				UploadedAt: time.Date(2023, time.Month(2), 10, 1, 10, 30, 0, timezone2),
				User:       "395fd5f4-964d-4135-9a55-fbf91c4a163b",
			},
		}, nil).
		AnyTimes()

	m.EXPECT().
		UserOrders(
			gomock.Any(),
			"395fd5f4-964d-4135-9a55-fbf91c4a1613",
		).
		Return([]*entity.Order{}, nil).
		AnyTimes()

	cfg := config.Config{
		SrvAddr:   "8000",
		SecretKey: "yvdUuY)HSX}?&b",
	}

	orderUseCase := order.New(m)
	userUseCase := user.New(nil, cfg.SecretKey)

	jwtToken, _ := auth.GenerateToken(
		cfg.SecretKey,
		&entity.User{ID: "395fd5f4-964d-4135-9a55-fbf91c4a163b"},
	)

	jwtToken2, _ := auth.GenerateToken(
		cfg.SecretKey,
		&entity.User{ID: "395fd5f4-964d-4135-9a55-fbf91c4a1613"},
	)

	ts := httptest.NewServer(
		NewRouter(
			&cfg,
			NewHandler(&cfg, nil, nil, orderUseCase, userUseCase),
		),
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
				"[{\"number\":\"2774311589\",\"status\":\"NEW\",\"accrual\":0,\"uploaded_at\":\"2023-03-11T01:10:30-08:00\"},{\"number\":\"12345678903\",\"status\":\"PROCESSED\",\"accrual\":100,\"uploaded_at\":\"2023-02-21T01:10:30+00:00\"},{\"number\":\"252576137\",\"status\":\"INVALID\",\"accrual\":0,\"uploaded_at\":\"2023-02-10T01:10:30+05:00\"}]\n",
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
