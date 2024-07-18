package main

import (
	"net/http"

	"github.com/alecthomas/kingpin/v2"
	"github.com/lvoytek/discourse_client_go/pkg/discourse"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
)

func main() {
	var (
		discourseSiteURL = kingpin.Flag("discourse.site-url", "The URL of the Discourse site to collect metrics from.").Default("http://127.0.0.1:3000").String()
		listenAddress    = kingpin.Flag("web.listen-address", "The address on which to expose the web interface and generated Prometheus metrics.").Default(":10110").String()
		metricsEndpoint  = kingpin.Flag("web.telemetry-path", "Endpoint path where metrics are exposed.").Default("/metrics").String()
	)

	kingpin.Version(version.Print("discourse_exporter"))
	kingpin.Parse()

	discourseClient := discourse.NewAnonymousClient(*discourseSiteURL)

	http.Handle(*metricsEndpoint, promhttp.Handler())
	http.ListenAndServe(*listenAddress, nil)
}
