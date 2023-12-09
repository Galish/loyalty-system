package router

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/Galish/loyalty-system/internal/auth"
	"github.com/Galish/loyalty-system/internal/config"
	"github.com/Galish/loyalty-system/internal/repository"
	"github.com/Galish/loyalty-system/internal/repository/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlerRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockUserRepository(ctrl)

	m.EXPECT().
		Create(
			gomock.Any(),
			"username",
			gomock.Any(),
		).
		DoAndReturn(func(ctx context.Context, login, password string) (*repository.User, error) {
			return &repository.User{
				Login:    login,
				Password: password,
			}, nil
		}).
		AnyTimes()

	m.EXPECT().
		Create(
			gomock.Any(),
			"user",
			gomock.Any(),
		).
		Return(nil, repository.ErrUserConflict).
		AnyTimes()

	cfg := config.Config{SrvAddr: "8000"}
	authService := auth.NewService(m)

	ts := httptest.NewServer(
		New(&cfg, authService),
	)
	defer ts.Close()

	type request struct {
		path   string
		method string
		body   *auth.Credentials
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
				"/api/user/registration",
				"POST",
				&auth.Credentials{},
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
				"/api/user/register",
				"GET",
				&auth.Credentials{},
			},
			&want{
				http.StatusMethodNotAllowed,
				"",
				"",
			},
		},
		{
			"empty request body",
			&request{
				"/api/user/register",
				"POST",
				&auth.Credentials{},
			},
			&want{
				http.StatusBadRequest,
				"missing login or password\n",
				"text/plain; charset=utf-8",
			},
		},
		{
			"missing password",
			&request{
				"/api/user/register",
				"POST",
				&auth.Credentials{
					Login: "username",
				},
			},
			&want{
				http.StatusBadRequest,
				"missing login or password\n",
				"text/plain; charset=utf-8",
			},
		},
		{
			"missing login",
			&request{
				"/api/user/register",
				"POST",
				&auth.Credentials{
					Password: "123456",
				},
			},
			&want{
				http.StatusBadRequest,
				"missing login or password\n",
				"text/plain; charset=utf-8",
			},
		},
		{
			"successful registration",
			&request{
				"/api/user/register",
				"POST",
				&auth.Credentials{
					Login:    "username",
					Password: "123456",
				},
			},
			&want{
				http.StatusOK,
				"",
				"",
			},
		},
		{
			"user already exists",
			&request{
				"/api/user/register",
				"POST",
				&auth.Credentials{
					Login:    "user",
					Password: "123456",
				},
			},
			&want{
				http.StatusConflict,
				"user already exists\n",
				"text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, err := json.Marshal(tt.req.body)
			require.NoError(t, err)

			req, err := http.NewRequest(
				tt.req.method,
				ts.URL+tt.req.path,
				bytes.NewBuffer(reqBody),
			)
			require.NoError(t, err)

			// Disable compression
			// req.Header.Set("Accept-Encoding", "identity")

			client := &http.Client{}
			resp, err := client.Do(req)
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, resp.StatusCode)

			raw, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, resp.Header.Get("Content-Type"), tt.want.contentType)
			assert.Equal(t, tt.want.body, string(raw))

			if resp.StatusCode == 200 {
				var authCookie *http.Cookie

				for _, c := range resp.Cookies() {
					if c.Name == auth.AuthCookieName {
						authCookie = c
						break
					}
				}

				assert.Regexp(
					t,
					regexp.MustCompile("^[0-9A-Za-z-_.]{10,}$"),
					authCookie.Value,
				)
			}

			err = resp.Body.Close()
			require.NoError(t, err)
		})
	}
}

func TestHandlerLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockUserRepository(ctrl)

	m.EXPECT().
		GetByLogin(
			gomock.Any(),
			"username",
		).
		DoAndReturn(func(ctx context.Context, login string) (*repository.User, error) {
			return &repository.User{
				Login:    login,
				Password: "$2a$10$5B2fWcB3sw2ONZ25klPPMe688GlUjjAUHRV.HnCd7xxG.KFX3CwBi",
			}, nil
		}).
		AnyTimes()

	m.EXPECT().
		GetByLogin(
			gomock.Any(),
			"user",
		).
		Return(nil, repository.ErrUserNotFound).
		AnyTimes()

	cfg := config.Config{SrvAddr: "8000"}
	authService := auth.NewService(m)

	ts := httptest.NewServer(
		New(&cfg, authService),
	)
	defer ts.Close()

	type request struct {
		path   string
		method string
		body   *auth.Credentials
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
				"/api/user/auth",
				"POST",
				&auth.Credentials{},
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
				"/api/user/login",
				"GET",
				&auth.Credentials{},
			},
			&want{
				http.StatusMethodNotAllowed,
				"",
				"",
			},
		},
		{
			"empty request body",
			&request{
				"/api/user/login",
				"POST",
				&auth.Credentials{},
			},
			&want{
				http.StatusBadRequest,
				"missing login or password\n",
				"text/plain; charset=utf-8",
			},
		},
		{
			"missing password",
			&request{
				"/api/user/login",
				"POST",
				&auth.Credentials{
					Login: "username",
				},
			},
			&want{
				http.StatusBadRequest,
				"missing login or password\n",
				"text/plain; charset=utf-8",
			},
		},
		{
			"missing login",
			&request{
				"/api/user/login",
				"POST",
				&auth.Credentials{
					Password: "123456",
				},
			},
			&want{
				http.StatusBadRequest,
				"missing login or password\n",
				"text/plain; charset=utf-8",
			},
		},
		{
			"successful login",
			&request{
				"/api/user/login",
				"POST",
				&auth.Credentials{
					Login:    "username",
					Password: "123456",
				},
			},
			&want{
				http.StatusOK,
				"",
				"",
			},
		},
		{
			"user not found",
			&request{
				"/api/user/login",
				"POST",
				&auth.Credentials{
					Login:    "user",
					Password: "123456",
				},
			},
			&want{
				http.StatusUnauthorized,
				"user not found\n",
				"text/plain; charset=utf-8",
			},
		},
		{
			"incorrect login/password pair",
			&request{
				"/api/user/login",
				"POST",
				&auth.Credentials{
					Login:    "username",
					Password: "12345678",
				},
			},
			&want{
				http.StatusUnauthorized,
				"incorrect login/password pair\n",
				"text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, err := json.Marshal(tt.req.body)
			require.NoError(t, err)

			req, err := http.NewRequest(
				tt.req.method,
				ts.URL+tt.req.path,
				bytes.NewBuffer(reqBody),
			)
			require.NoError(t, err)

			// Disable compression
			// req.Header.Set("Accept-Encoding", "identity")

			client := &http.Client{}
			resp, err := client.Do(req)
			require.NoError(t, err)

			assert.Equal(t, tt.want.statusCode, resp.StatusCode)

			raw, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, resp.Header.Get("Content-Type"), tt.want.contentType)
			assert.Equal(t, tt.want.body, string(raw))

			if resp.StatusCode == 200 {
				var authCookie *http.Cookie

				for _, c := range resp.Cookies() {
					if c.Name == auth.AuthCookieName {
						authCookie = c
						break
					}
				}

				assert.Regexp(
					t,
					regexp.MustCompile("^[0-9A-Za-z-_.]{10,}$"),
					authCookie.Value,
				)
			}

			err = resp.Body.Close()
			require.NoError(t, err)
		})
	}
}
