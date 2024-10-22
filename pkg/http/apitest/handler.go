package apitest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHandler is a helper function to test http.Handler implementations.
func TestHandler(t *testing.T, r *http.Request, wantCode int, wantBody string, h http.Handler) {
	t.Helper()

	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	res := w.Result()
	assert.Equal(t, wantCode, res.StatusCode)

	if wantBody != "" {
		switch res.Header.Get("Content-Type") {
		case "application/json":
			require.JSONEq(t, wantBody, w.Body.String())
		default:
			require.Equal(t, wantBody, w.Body.String())
		}
	} else {
		require.Empty(t, w.Body.String())
	}
}
