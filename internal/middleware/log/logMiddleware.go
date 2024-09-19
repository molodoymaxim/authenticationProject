package log

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var log = logrus.New()

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := middleware.GetReqID(r.Context())
		start := time.Now()

		log.WithFields(logrus.Fields{
			"requestID": reqID,
			"method":    r.Method,
			"path":      r.URL.Path,
		}).Info("Started request")

		next.ServeHTTP(w, r)

		log.WithFields(logrus.Fields{
			"requestID": reqID,
			"method":    r.Method,
			"path":      r.URL.Path,
			"duration":  time.Since(start),
		}).Info("Completed request")
	})
}
