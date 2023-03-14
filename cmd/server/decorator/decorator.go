package decorator

import "net/http"

type HandleWithError func(w http.ResponseWriter, r *http.Request) error

type Decorator func(HandleWithError) HandleWithError

func HTTP(handle HandleWithError, ds ...Decorator) http.Handler {
	for _, d := range ds {
		handle = d(handle)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = handle(w, r)
	})
}
