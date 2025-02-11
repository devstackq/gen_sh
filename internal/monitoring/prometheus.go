package monitoring

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	VideosGenerated = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "videos_generated_total",
			Help: "Total number of generated videos",
		},
	)
)

func InitMonitoring() {
	prometheus.MustRegister(VideosGenerated)

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Println("📊 Запуск мониторинга на :9090")
		http.ListenAndServe(":9090", nil)
	}()
}
