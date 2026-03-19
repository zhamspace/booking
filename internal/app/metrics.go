package app

import (
	"log/slog"
	"os"

	"github.com/zhamspace/booking/internal/config"
	"github.com/zhamspace/booking/internal/constant"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	metricRequestCounter   *prometheus.CounterVec
	metricResponseDuration *prometheus.HistogramVec
)

func init() {
	if !config.Conf.WithMetrics {
		return
	}

	slog.New(slog.NewJSONHandler(os.Stdout, nil)).Info("metrics enabled")

	metricRequestCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: config.Conf.Namespace,
		Name:      constant.ServiceName + "_request_count",
	}, []string{
		"protocol",
		"method",
		"status",
	})

	metricResponseDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: config.Conf.Namespace,
		Name:      constant.ServiceName + "_response_duration_seconds",
	}, []string{
		"protocol",
		"method",
		"status",
	})
}
