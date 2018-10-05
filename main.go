package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/silva-willian/jenkins_exporter/collector"
)

func main() {
	jenkinsCollector := collector.JenkinsMetricsCollector()
	prometheus.MustRegister(jenkinsCollector)
	http.Handle("/metrics", promhttp.Handler())

	var port = os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	fmt.Println("Server running in port:", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
