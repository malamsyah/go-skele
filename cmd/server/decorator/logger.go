package decorator

import (
	"net/http"

	"github.com/felixge/httpsnoop"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// Logger will log every request. Requires to be put into mux router, otherwise it will panic
// nolint
func Logger() Decorator {
	return func(next HandleWithError) HandleWithError {
		return func(w http.ResponseWriter, r *http.Request) error {
			logger, _ := zap.NewProduction()
			defer logger.Sync()
			sugar := logger.Sugar()

			var err error
			metric := httpsnoop.CaptureMetricsFn(w, func(ww http.ResponseWriter) {
				err = next(ww, r)
			})

			path, _ := mux.CurrentRoute(r).GetPathTemplate()
			logger.With(zap.Field{
				Key:       "",
				Type:      0,
				Integer:   0,
				String:    "",
				Interface: nil,
			})
			sugar.With("method", r.Method,
				"path", path,
				"code", metric.Code,
				"duration", metric.Duration.Milliseconds(),
				"params", mux.Vars(r),
				"query_string", r.URL.Query())

			return err
		}
	}
}
