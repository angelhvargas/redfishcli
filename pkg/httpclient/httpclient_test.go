package httpclient

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/angelhvargas/redfishcli/pkg/logger"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestDoRequest(t *testing.T) {
	// Set up a logger with a test hook
	var log = logrus.New()
	hook := test.NewLocal(log)

	// Replace the logger with the test logger
	logger.Log = log

	t.Run("Successful request", func(t *testing.T) {
		// Reset the hook entries before each test
		hook.Reset()

		// Create a mock server that returns a successful response
		mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "Basic dXNlcjpwYXNz", r.Header.Get("Authorization"))
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("success"))
		}))
		defer mockServer.Close()

		// Configure the HTTP client to trust the mock server's certificate
		config := Config{
			Timeout:       90 * time.Second,
			SkipTLSVerify: true,
		}

		// Perform the request
		body, err := DoRequest(mockServer.URL, "user", "pass", config)
		assert.NoError(t, err)
		assert.Equal(t, "success", string(body))

		// Check log entries
		assert.Len(t, hook.Entries, 1)
		assert.Equal(t, logrus.InfoLevel, hook.LastEntry().Level)
		assert.Contains(t, hook.LastEntry().Message, "Doing http request to")
	})

	t.Run("Request with invalid URL", func(t *testing.T) {
		// Reset the hook entries before each test
		hook.Reset()

		config := Config{
			Timeout:       90 * time.Second,
			SkipTLSVerify: true,
		}

		_, err := DoRequest(":", "user", "pass", config)
		assert.Error(t, err)

		// Check log entries
		assert.Len(t, hook.Entries, 1)
		assert.Equal(t, logrus.InfoLevel, hook.LastEntry().Level)
		assert.Contains(t, hook.LastEntry().Message, "Doing http request to")
	})

	t.Run("Request with server error", func(t *testing.T) {
		// Reset the hook entries before each test
		hook.Reset()

		// Create a mock server that returns an error response
		mockServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer mockServer.Close()

		// Configure the HTTP client to trust the mock server's certificate
		config := Config{
			Timeout:       90 * time.Second,
			SkipTLSVerify: true,
		}

		// Perform the request
		_, err := DoRequest(mockServer.URL, "user", "pass", config)
		assert.Error(t, err)

		// Check log entries
		assert.Len(t, hook.Entries, 2)
		assert.Equal(t, logrus.InfoLevel, hook.Entries[0].Level)
		assert.Contains(t, hook.Entries[0].Message, "Doing http request to")
		assert.Equal(t, logrus.ErrorLevel, hook.Entries[1].Level)
		assert.Contains(t, hook.Entries[1].Message, "Error")
	})
}
