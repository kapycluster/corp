package middleware

import (
	"context"
	"net/http"

	"github.com/go-chi/httplog/v2"
	"github.com/kapycluster/corpy/log"
)

func RequestLogger(ctx context.Context) func(next http.Handler) http.Handler {
	httpLogger := &httplog.Logger{
		Logger: log.FromContext(ctx),
		Options: httplog.Options{
			Concise: true,
		},
	}
	return httplog.Handler(httpLogger, []string{})
}
