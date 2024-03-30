package middleware

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"firebase.google.com/go/v4/auth"
)

func TestAuthenticationMiddleware(t *testing.T) {
	type args struct {
		verifyToken    func(context.Context, string) (*auth.Token, error)
		allowedMethods []string
		getenv         func(string) string
	}
	type req struct {
		method string
		url    string
		body   io.Reader
		token  string
	}
	tests := []struct {
		name string
		args args
		req  req
		want int
	}{
		{
			name: "test1",
			args: args{
				verifyToken: func(ctx context.Context, token string) (*auth.Token, error) {
					return nil, nil
				},
				allowedMethods: []string{http.MethodGet},
				getenv:         func(s string) string { return "" },
			},
			req: req{
				method: http.MethodDelete,
				url:    "/test",
				body:   nil,
			},
			want: http.StatusUnauthorized,
		},
		{
			name: "test2",
			args: args{
				verifyToken: func(ctx context.Context, token string) (*auth.Token, error) {
					return nil, nil
				},
				allowedMethods: []string{http.MethodGet},
				getenv:         func(s string) string { return "" },
			},
			req: req{
				method: http.MethodGet,
				url:    "/test",
				body:   nil,
			},
			want: http.StatusOK,
		},
		{
			name: "test3",
			args: args{
				verifyToken: func(ctx context.Context, token string) (*auth.Token, error) {
					return &auth.Token{}, nil
				},
				allowedMethods: []string{http.MethodDelete},
				getenv:         func(s string) string { return "" },
			},
			req: req{
				method: http.MethodGet,
				url:    "/test",
				body:   nil,
				token:  "123456",
			},
			want: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.req.method, tt.req.url, tt.req.body)
			if tt.req.token != "" {
				r.Header.Add("x-api-token", tt.req.token)
			}
			w := httptest.NewRecorder()

			sut := AuthenticationMiddleware(tt.args.verifyToken, tt.args.allowedMethods, tt.args.getenv)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

			sut.ServeHTTP(w, r)

			if !strings.HasPrefix(w.Result().Status, fmt.Sprintf("%d", tt.want)) {
				t.Errorf("middleware: AuthenticationMiddleware - got: %s want: %d", w.Result().Status, tt.want)
			}
		})
	}
}
