package httpclient

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/angelhvargas/redfishcli/pkg/logger"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Custom hook for logrus to capture logs for assertions
type TestLogHook struct {
	Entries []*logrus.Entry
}

func (hook *TestLogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *TestLogHook) Fire(entry *logrus.Entry) error {
	hook.Entries = append(hook.Entries, entry)
	return nil
}

// Setup a test logger
func SetupTestLogger() (*logrus.Logger, *TestLogHook) {
	logger := logrus.New()
	hook := &TestLogHook{}
	logger.AddHook(hook)
	return logger, hook
}

func TestDoRequest(t *testing.T) {
	// Setup the test logger
	testLogger, testLogHook := SetupTestLogger()
	originalLogger := logger.Log
	logger.Log = testLogger
	defer func() { logger.Log = originalLogger }()

	t.Run("Successful_request", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusOK)
			rw.Write([]byte("OK"))
		}))
		defer server.Close()

		config := DefaultConfig()
		body, err := DoRequest(server.URL, "user", "pass", config)
		require.NoError(t, err)
		assert.Equal(t, "OK", string(body))

		require.Len(t, testLogHook.Entries, 2)
		assert.Contains(t, testLogHook.Entries[0].Message, "Doing http request to")
		assert.Contains(t, testLogHook.Entries[1].Message, "200")
	})

	t.Run("Request_with_invalid_URL", func(t *testing.T) {
		config := DefaultConfig()
		_, err := DoRequest(":", "user", "pass", config)
		require.Error(t, err)

		require.Len(t, testLogHook.Entries, 3)
		assert.Contains(t, testLogHook.Entries[2].Message, "Doing http request to")
	})

	t.Run("Request_with_server_error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		config := DefaultConfig()
		_, err := DoRequest(server.URL, "user", "pass", config)
		require.Error(t, err)

		require.Len(t, testLogHook.Entries, 6)
		assert.Contains(t, testLogHook.Entries[3].Message, "Doing http request to")
		assert.Contains(t, testLogHook.Entries[4].Message, "500")
		assert.Contains(t, testLogHook.Entries[5].Message, "Error: server returned status code 500")
	})
}
