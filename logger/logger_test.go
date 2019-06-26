package logger

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	var buf bytes.Buffer

	l := New(&buf, "somelogger", "info", "key1", "value1")

	t.Run("info", func(t *testing.T) {
		defer buf.Reset()

		l.Info("foo", "key2", "value2")

		require.Contains(t, buf.String(), "key1=value1")
		require.Contains(t, buf.String(), "key2=value2")
		require.Contains(t, buf.String(), "caller=github.com/graymeta/gmkit/logger/logger_test.go")
		require.Contains(t, buf.String(), "ts="+fmt.Sprintf("%d", time.Now().Year()))
		require.Contains(t, buf.String(), "src=somelogger")
		require.Contains(t, buf.String(), "level=info")
		require.Contains(t, buf.String(), "msg=foo")

	})

	t.Run("check-level-filtering", func(t *testing.T) {
		defer buf.Reset()

		l.Debug("debug_message", "keyDebug", "valueDebug")
		require.NotContains(t, buf.String(), "debug_message")
	})

	t.Run("sub-logger", func(t *testing.T) {
		defer buf.Reset()

		sub := l.New("sublogger")
		sub.Info("sub", "key3", "value3")

		require.Contains(t, buf.String(), "key1=value1")
		require.Contains(t, buf.String(), "key3=value3")
		require.Contains(t, buf.String(), "src=somelogger.sublogger")
	})

	t.Run("sub-log-with-vals", func(t *testing.T) {
		defer buf.Reset()

		sub := l.New("sublogger2").With("key4", "value4")
		sub.Info("sub", "key5", "value5")

		require.Contains(t, buf.String(), "key1=value1")
		require.Contains(t, buf.String(), "key4=value4")
		require.Contains(t, buf.String(), "key5=value5")
		require.Contains(t, buf.String(), "src=somelogger.sublogger2")
	})
}
