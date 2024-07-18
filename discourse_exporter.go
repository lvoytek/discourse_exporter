package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/lvoytek/discourse_client_go/pkg/discourse"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
)

func main() {
	var (
		discourseSiteURL       = kingpin.Flag("discourse.site-url", "The URL of the Discourse site to collect metrics from.").Default("http://127.0.0.1:3000").String()
		discourseCategoryList  = kingpin.Flag("discourse.limit-categories", "Comma separated list of category slugs to limit metrics to. All are enabled by default.").Default("").String()
		dataCollectionInterval = kingpin.Flag("discourse.collection-interval", "Time in seconds to wait before collecting new data from the Discourse site.").Default("3600").Int()
		listenAddress          = kingpin.Flag("web.listen-address", "The address on which to expose the web interface and generated Prometheus metrics.").Default(":10110").String()
		metricsEndpoint        = kingpin.Flag("web.telemetry-path", "Endpoint path where metrics are exposed.").Default("/metrics").String()
	)

	kingpin.Version(version.Print("discourse_exporter"))
	kingpin.Parse()

	prometheus.MustRegister(version.NewCollector("discourse_exporter"))

	discourseClient := discourse.NewAnonymousClient(*discourseSiteURL)

	go IntervalCollect(discourseClient, strings.Split(strings.TrimSpace(*discourseCategoryList), ","), time.Duration(*dataCollectionInterval)*time.Second)

	http.Handle(*metricsEndpoint, promhttp.Handler())
	http.ListenAndServe(*listenAddress, nil)
}
