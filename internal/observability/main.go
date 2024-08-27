package observability

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"

	dd_logrus "gopkg.in/DataDog/dd-trace-go.v1/contrib/sirupsen/logrus"
)

var log = logrus.New()

func Init() {
	log.SetFormatter(&logrus.JSONFormatter{})
	log.Info("Logrus set to JSON formatter")

	// Output to stdout instead of the default stderr
	log.SetOutput(os.Stdout)
	log.Info("Logrus set to output to stdout")

	// Only log the info severity or above
	log.SetLevel(logrus.InfoLevel)

	// Add Datadog context log hook
	log.AddHook(&dd_logrus.DDContextLogHook{})
}

func InfoWithContext(ctx context.Context, msg string) {
    log.WithContext(ctx).Info(msg)
}

func ErrorWithContext(ctx context.Context, msg string) {
    log.WithContext(ctx).Error(msg)
}
