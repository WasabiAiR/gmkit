package logger

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/cheekybits/is"
)

func TestLogger(t *testing.T) {
	var buf bytes.Buffer

	l := New(&buf, "somelogger", "info", "key1", "value1")

	t.Run("info", func(t *testing.T) {
		is := is.New(t)
		defer buf.Reset()

		l.Info("foo", "key2", "value2")

		is.True(strings.Contains(buf.String(), "key1=value1"))
		is.True(strings.Contains(buf.String(), "key2=value2"))
		is.True(strings.Contains(buf.String(), "caller=github.com/graymeta/saas/internal/logger/logger_test.go"))
		is.True(strings.Contains(buf.String(), "ts="+fmt.Sprintf("%d", time.Now().Year())))
		is.True(strings.Contains(buf.String(), "src=somelogger"))
		is.True(strings.Contains(buf.String(), "level=info"))
		is.True(strings.Contains(buf.String(), "msg=foo"))

	})

	t.Run("check-level-filtering", func(t *testing.T) {
		is := is.New(t)
		defer buf.Reset()

		l.Debug("debug_message", "keyDebug", "valueDebug")
		is.False(strings.Contains(buf.String(), "debug_message"))
	})

	t.Run("sub-logger", func(t *testing.T) {
		is := is.New(t)
		defer buf.Reset()

		sub := l.New("sublogger")
		sub.Info("sub", "key3", "value3")

		is.True(strings.Contains(buf.String(), "key1=value1"))
		is.True(strings.Contains(buf.String(), "key3=value3"))
		is.True(strings.Contains(buf.String(), "src=somelogger.sublogger"))
	})

	t.Run("sub-log-with-vals", func(t *testing.T) {
		is := is.New(t)
		defer buf.Reset()

		sub := l.New("sublogger2").With("key4", "value4")
		sub.Info("sub", "key5", "value5")

		is.True(strings.Contains(buf.String(), "key1=value1"))
		is.True(strings.Contains(buf.String(), "key4=value4"))
		is.True(strings.Contains(buf.String(), "key5=value5"))
		is.True(strings.Contains(buf.String(), "src=somelogger.sublogger2"))
	})
}
