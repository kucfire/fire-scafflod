package monitor

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const ()

func PrometheusHTTPServer() http.Handler {
	return promhttp.Handler()
}
