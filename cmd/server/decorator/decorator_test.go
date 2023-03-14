package decorator_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/malamsyah/go-skele/cmd/server/decorator"
	"github.com/stretchr/testify/assert"
)

func TestHTTP(t *testing.T) {
	dummyHandle := func(w http.ResponseWriter, r *http.Request) error {
		_, _ = w.Write([]byte(`{"message": "ok"}`))

		return nil
	}

	dummyDecorator1 := func(next decorator.HandleWithError) decorator.HandleWithError {
		return func(w http.ResponseWriter, r *http.Request) error {
			return next(w, r)
		}
	}

	dummyDecorator2 := func(next decorator.HandleWithError) decorator.HandleWithError {
		return func(w http.ResponseWriter, r *http.Request) error {
			return next(w, r)
		}
	}

	assert.NotPanics(t, func() {
		handle := decorator.HTTP(dummyHandle, dummyDecorator1, dummyDecorator2)

		w := httptest.NewRecorder()
		r, _ := http.NewRequest(http.MethodGet, "/", http.NoBody)
		handle.ServeHTTP(w, r)

		assert.Equal(t, `{"message": "ok"}`, strings.TrimSpace(w.Body.String()))
		assert.NotNil(t, handle)
	})
}
