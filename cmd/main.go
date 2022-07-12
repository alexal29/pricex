package main

import (
	"flag"
	"net/http"

	"github.com/alexal29/pricex/pkg/handlers"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	log "github.com/sirupsen/logrus"
)

var (
	listenAddress = flag.String("web.listen-address", "0.0.0.0:9104",
		"listen address")
	metricsPath = flag.String("web.metric-path", "/metrics",
		"Path under which to expose metrics")
)

func main() {
	// =====================
	// Get OS parameter
	// =====================
	flag.Parse()

	// ========================
	// Regist handler
	// ========================
	log.Infof("Regist version collector pricex")
	prometheus.Register(version.NewCollector("pricex"))
	prometheus.Register(handlers.NewPriceCollector())

	// Regist http handler
	http.HandleFunc(*metricsPath, func(w http.ResponseWriter, r *http.Request) {
		h := promhttp.HandlerFor(prometheus.Gatherers{
			prometheus.DefaultGatherer,
		}, promhttp.HandlerOpts{})
		h.ServeHTTP(w, r)
	})

	// start server
	log.Infof("Starting http server - %s", *listenAddress)
	if err := http.ListenAndServe(*listenAddress, nil); err != nil {
		log.Errorf("Failed to start http server: %s", err)
	}
}
