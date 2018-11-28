package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dasa-health/jenkins_exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func handlerJobs(w http.ResponseWriter, r *http.Request) {
	registry := prometheus.NewRegistry()
	collector := collector.JobsMetricsCollector()
	registry.MustRegister(collector)
	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}

func handlerMetrics(w http.ResponseWriter, r *http.Request) {
	registry := prometheus.NewRegistry()
	collector := collector.JenkinsMetricsCollector(r.URL.Query().Get("jobName"))
	registry.MustRegister(collector)
	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}

func main() {
	http.HandleFunc("/metrics", handlerMetrics)
	http.HandleFunc("/jobs", handlerJobs)

	var port = os.Getenv("port")

	if port == "" {
		port = "3000"
	}

	fmt.Println("Server running in port:", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
