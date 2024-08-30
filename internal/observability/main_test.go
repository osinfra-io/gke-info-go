package observability

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type TestHook struct {
    Entries []*logrus.Entry
}

func (hook *TestHook) Levels() []logrus.Level {
    return logrus.AllLevels
}

func (hook *TestHook) Fire(entry *logrus.Entry) error {
    hook.Entries = append(hook.Entries, entry)
    return nil
}

func TestInit(t *testing.T) {
    // Add a test hook to capture log entries
    hook := &TestHook{}
    log.AddHook(hook)

    // Initialize the logger
    Init()

    // Check if the logger is set to JSON formatter
    _, ok := log.Formatter.(*logrus.JSONFormatter)
    assert.True(t, ok, "Expected JSONFormatter")

    // Check if the logger level is set to InfoLevel
    assert.Equal(t, logrus.InfoLevel, log.Level, "Expected log level to be InfoLevel")

    // Check if the initialization log messages are present
    assert.Len(t, hook.Entries, 2, "Expected two log entries during initialization")
    assert.Contains(t, hook.Entries[0].Message, "Logrus set to JSON formatter", "Expected first log message to contain 'Logrus set to JSON formatter'")
    assert.Contains(t, hook.Entries[1].Message, "Logrus set to output to stdout", "Expected second log message to contain 'Logrus set to output to stdout'")
}

func TestInfoWithContext(t *testing.T) {
    // Add a test hook to capture log entries
    hook := &TestHook{}
    log.AddHook(hook)

    // Create a context
    ctx := context.Background()

    // Log an info message with context
    InfoWithContext(ctx, "test info message")

    // Check if the message is logged
    assert.Len(t, hook.Entries, 1, "Expected one log entry")
    assert.Contains(t, hook.Entries[0].Message, "test info message", "Expected log message to contain 'test info message'")
}

func TestErrorWithContext(t *testing.T) {
    // Add a test hook to capture log entries
    hook := &TestHook{}
    log.AddHook(hook)

    // Create a context
    ctx := context.Background()

    // Log an error message with context
    ErrorWithContext(ctx, "test error message")

    // Check if the message is logged
    assert.Len(t, hook.Entries, 1, "Expected one log entry")
    assert.Contains(t, hook.Entries[0].Message, "test error message", "Expected log message to contain 'test error message'")
}
