package pgmfxmetrics

import (
	"fmt"
	"time"

	"github.com/myfantasy/mfctx"
	"github.com/myfantasy/pgmfx"
	"github.com/prometheus/client_golang/prometheus"
)

type MetricsCommon struct {
	Request *prometheus.CounterVec

	GetPoolTotal     *prometheus.CounterVec
	GetPoolTimeHist  *prometheus.HistogramVec
	GetPoolTimeTotal *prometheus.CounterVec

	DoRequestTotal     *prometheus.CounterVec
	DoRequestTimeHist  *prometheus.HistogramVec
	DoRequestTimeTotal *prometheus.CounterVec

	DoReadResponceTotal     *prometheus.CounterVec
	DoReadResponceTimeHist  *prometheus.HistogramVec
	DoReadResponceTimeTotal *prometheus.CounterVec

	ResponceTotal     *prometheus.CounterVec
	ResponceTimeHist  *prometheus.HistogramVec
	ResponceTimeTotal *prometheus.CounterVec
}

var _ pgmfx.MetricsProvider = &MetricsCommon{}

func NewMetricsCommon() *MetricsCommon {
	constLabels := prometheus.Labels{
		"version":  mfctx.GetAppVersion(),
		"app_name": mfctx.GetAppName(),
		"app_id":   mfctx.GetAppID(),
	}

	return &MetricsCommon{
		Request: prometheus.NewCounterVec(
			// nolint:promlinter
			prometheus.CounterOpts{
				Name:        "pg_request_start",
				Help:        "Total amount of runnings request",
				ConstLabels: constLabels,
			}, []string{"pool_name", "sql_name"},
		),

		GetPoolTotal: prometheus.NewCounterVec(
			// nolint:promlinter
			prometheus.CounterOpts{
				Name:        "pg_get_pool_finish",
				Help:        "How many get pool runnings finish, partitioned",
				ConstLabels: constLabels,
			}, []string{
				"pool_name", "sql_name", "status_code", "alarm",
			},
		),
		GetPoolTimeHist: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:        "pg_get_pool_finish_hist",
				Help:        "Total amount of time spent on the runnings finish",
				ConstLabels: constLabels,
				Buckets:     []float64{5, 10, 20, 50, 100, 200, 500, 1000, 2000},
			}, []string{
				"pool_name", "sql_name", "status_code", "alarm",
			},
		),
		GetPoolTimeTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "pg_get_pool_finish_time_total",
				Help:        "Total amount of time spent on the runnings finish",
				ConstLabels: constLabels,
			}, []string{
				"pool_name", "sql_name", "status_code", "alarm",
			},
		),

		DoRequestTotal: prometheus.NewCounterVec(
			// nolint:promlinter
			prometheus.CounterOpts{
				Name:        "pg_do_request_finish",
				Help:        "How many get pool runnings finish, partitioned",
				ConstLabels: constLabels,
			}, []string{
				"pool_name", "sql_name", "status_code", "alarm",
			},
		),
		DoRequestTimeHist: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:        "pg_do_request_finish_hist",
				Help:        "Total amount of time spent on the runnings finish",
				ConstLabels: constLabels,
				Buckets:     []float64{5, 10, 20, 50, 100, 200, 500, 1000, 2000},
			}, []string{
				"pool_name", "sql_name", "status_code", "alarm",
			},
		),
		DoRequestTimeTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "pg_do_request_finish_time_total",
				Help:        "Total amount of time spent on the runnings finish",
				ConstLabels: constLabels,
			}, []string{
				"pool_name", "sql_name", "status_code", "alarm",
			},
		),

		DoReadResponceTotal: prometheus.NewCounterVec(
			// nolint:promlinter
			prometheus.CounterOpts{
				Name:        "pg_do_read_finish",
				Help:        "How many get pool runnings finish, partitioned",
				ConstLabels: constLabels,
			}, []string{
				"pool_name", "sql_name", "status_code", "alarm",
			},
		),
		DoReadResponceTimeHist: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:        "pg_do_read_finish_hist",
				Help:        "Total amount of time spent on the runnings finish",
				ConstLabels: constLabels,
				Buckets:     []float64{5, 10, 20, 50, 100, 200, 500, 1000, 2000},
			}, []string{
				"pool_name", "sql_name", "status_code", "alarm",
			},
		),
		DoReadResponceTimeTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "pg_do_read_finish_time_total",
				Help:        "Total amount of time spent on the runnings finish",
				ConstLabels: constLabels,
			}, []string{
				"pool_name", "sql_name", "status_code", "alarm",
			},
		),

		ResponceTotal: prometheus.NewCounterVec(
			// nolint:promlinter
			prometheus.CounterOpts{
				Name:        "pg_request_finish",
				Help:        "How many requests runnings finish, partitioned",
				ConstLabels: constLabels,
			}, []string{
				"pool_name", "sql_name", "status_code", "alarm",
			},
		),
		ResponceTimeHist: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:        "pg_request_finish_hist",
				Help:        "Total amount of time spent on the runnings finish",
				ConstLabels: constLabels,
				Buckets:     []float64{5, 10, 20, 50, 100, 200, 500, 1000, 2000},
			}, []string{
				"pool_name", "sql_name", "status_code", "alarm",
			},
		),
		ResponceTimeTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name:        "pg_request_finish_time_total",
				Help:        "Total amount of time spent on the runnings finish",
				ConstLabels: constLabels,
			}, []string{
				"pool_name", "sql_name", "status_code", "alarm",
			},
		),
	}
}

func (mc *MetricsCommon) AutoRegister() *MetricsCommon {
	mc.MustRegister(prometheus.DefaultRegisterer)

	return mc
}

func (mc *MetricsCommon) MustRegister(registerer prometheus.Registerer) *MetricsCommon {
	registerer.MustRegister(
		mc.Request,

		mc.GetPoolTotal,
		mc.GetPoolTimeHist,
		mc.GetPoolTimeTotal,

		mc.DoRequestTotal,
		mc.DoRequestTimeHist,
		mc.DoRequestTimeTotal,

		mc.DoReadResponceTotal,
		mc.DoReadResponceTimeHist,
		mc.DoReadResponceTimeTotal,

		mc.ResponceTotal,
		mc.ResponceTimeHist,
		mc.ResponceTimeTotal,
	)

	return mc
}

func (mc *MetricsCommon) WriteMetricRequest(poolName, sqlName string) {
	if mc == nil {
		return
	}

	mc.Request.WithLabelValues(poolName, sqlName).Inc()
}

func (mc *MetricsCommon) WriteMetricResponce(mRequest time.Time, poolName, sqlName string, resultResult string) {
	if mc == nil {
		return
	}

	var alarm bool

	responseLabels := []string{poolName, sqlName, resultResult, fmt.Sprint(alarm)}
	diff := time.Since(mRequest).Milliseconds()

	mc.ResponceTotal.WithLabelValues(responseLabels...).Inc()
	mc.ResponceTimeHist.WithLabelValues(responseLabels...).Observe(float64(diff))
	mc.ResponceTimeTotal.WithLabelValues(responseLabels...).Add(float64(diff))
}

func (mc *MetricsCommon) WriteMetricGetPool(mRequest time.Time, poolName, sqlName string, resultResult string) {
	if mc == nil {
		return
	}

	var alarm bool

	responseLabels := []string{poolName, sqlName, resultResult, fmt.Sprint(alarm)}
	diff := time.Since(mRequest).Milliseconds()

	mc.GetPoolTotal.WithLabelValues(responseLabels...).Inc()
	mc.GetPoolTimeHist.WithLabelValues(responseLabels...).Observe(float64(diff))
	mc.GetPoolTimeTotal.WithLabelValues(responseLabels...).Add(float64(diff))
}

func (mc *MetricsCommon) WriteMetricDoRequest(mRequest time.Time, poolName, sqlName string, resultResult string) {
	if mc == nil {
		return
	}

	var alarm bool

	responseLabels := []string{poolName, sqlName, resultResult, fmt.Sprint(alarm)}
	diff := time.Since(mRequest).Milliseconds()

	mc.DoRequestTotal.WithLabelValues(responseLabels...).Inc()
	mc.DoRequestTimeHist.WithLabelValues(responseLabels...).Observe(float64(diff))
	mc.DoRequestTimeTotal.WithLabelValues(responseLabels...).Add(float64(diff))
}

func (mc *MetricsCommon) WriteMetricDoReadResponce(mRequest time.Time, poolName, sqlName string, resultResult string) {
	if mc == nil {
		return
	}

	var alarm bool

	responseLabels := []string{poolName, sqlName, resultResult, fmt.Sprint(alarm)}
	diff := time.Since(mRequest).Milliseconds()

	mc.DoReadResponceTotal.WithLabelValues(responseLabels...).Inc()
	mc.DoReadResponceTimeHist.WithLabelValues(responseLabels...).Observe(float64(diff))
	mc.DoReadResponceTimeTotal.WithLabelValues(responseLabels...).Add(float64(diff))
}
