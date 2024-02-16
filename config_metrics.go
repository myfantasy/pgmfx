package pgmfx

import "time"

const OKStatus = "OK"
const ErrorStatus = "ERR"
const ErrorPoolStatus = "ERR_P"
const ErrorQueryStatus = "ERR_Q"
const ErrorRowsReadStatus = "ERR_R"

var DefaultMetricsProvider MetricsProvider

type MetricsProvider interface {
	WriteMetricRequest(poolName, sqlName string)
	WriteMetricGetPool(mRequest time.Time, poolName, sqlName string, resultResult string)
	WriteMetricDoRequest(mRequest time.Time, poolName, sqlName string, resultResult string)
	WriteMetricDoReadResponce(mRequest time.Time, poolName, sqlName string, resultResult string)
	WriteMetricResponce(mRequest time.Time, poolName, sqlName string, resultResult string)
}

func WriteMetricRequest(mp MetricsProvider, poolName, sqlName string) {
	if mp == nil {
		mp = DefaultMetricsProvider
	}
	if mp != nil {
		mp.WriteMetricRequest(poolName, sqlName)
	}
}
func WriteMetricGetPool(mp MetricsProvider, mRequest time.Time, poolName, sqlName string, resultResult string) {
	if mp == nil {
		mp = DefaultMetricsProvider
	}
	if mp != nil {
		mp.WriteMetricGetPool(mRequest, poolName, sqlName, resultResult)
	}
}
func WriteMetricDoRequest(mp MetricsProvider, mRequest time.Time, poolName, sqlName string, resultResult string) {
	if mp == nil {
		mp = DefaultMetricsProvider
	}
	if mp != nil {
		mp.WriteMetricDoRequest(mRequest, poolName, sqlName, resultResult)
	}
}
func WriteMetricDoReadResponce(mp MetricsProvider, mRequest time.Time, poolName, sqlName string, resultResult string) {
	if mp == nil {
		mp = DefaultMetricsProvider
	}
	if mp != nil {
		mp.WriteMetricDoReadResponce(mRequest, poolName, sqlName, resultResult)
	}
}
func WriteMetricResponce(mp MetricsProvider, mRequest time.Time, poolName, sqlName string, resultResult string) {
	if mp == nil {
		mp = DefaultMetricsProvider
	}
	if mp != nil {
		mp.WriteMetricResponce(mRequest, poolName, sqlName, resultResult)
	}
}
