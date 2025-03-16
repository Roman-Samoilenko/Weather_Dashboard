package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type Middleware func(http.Handler) http.Handler

func CreateStack(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for _, mdlw := range xs {
			x := mdlw
			next = x(next)
		}
		return next
	}

}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		slog.Info(
			"request info",
			slog.Group("request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Duration("duration", time.Since(start)),
			),
		)
	})
}
