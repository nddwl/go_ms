package orm

import (
	"github.com/zeromicro/go-zero/core/metric"
)

const gormNameSpace = "gorm_client"

var (
	metricClientReqDur = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: gormNameSpace,
		Subsystem: "requests",
		Name:      "duration_ms",
		Help:      "gorm client requests duration(ms).",
		Labels:    []string{"table", "method"},
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000},
	})
	metricClientReqErrCount = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: gormNameSpace,
		Subsystem: "requests",
		Name:      "error_total",
		Help:      "gorm client requests error total.",
		Labels:    []string{"table", "method", "is_error"},
	})
)
