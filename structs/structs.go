package structs

import (
	"github.com/prometheus/client_golang/prometheus"
)

// JenkinsMetrics is the struct that maps the entities to the Prometheus client
type JenkinsMetrics struct {
	job   *prometheus.Desc
	stage *prometheus.Desc
}

// JenkinsStages is the struct that does the mapping with the stages of the Jenkis
type JenkinsStages struct {
	Name                string  `json:"name"`
	ID                  string  `json:"id"`
	Status              string  `json:"status"`
	StartTimeMillis     float64 `json:"startTimeMillis"`
	DurationMillis      float64 `json:"durationMillis"`
	PauseDurationMillis float64 `json:"pauseDurationMillis"`
	ExecNode            string  `json:"execNode"`
}

// Jenkins Job is the struct that represents the jobs of Jenkins
type JenkinsJob struct {
	Name                string          `json:"name"`
	ID                  string          `json:"id"`
	Status              string          `json:"status"`
	StartTimeMillis     float64         `json:"startTimeMillis"`
	EndTimeMillis       float64         `json:"endTimeMillis"`
	DurationMillis      float64         `json:"durationMillis"`
	QueueDurationMillis float64         `json:"queueDurationMillis"`
	PauseDurationMillis float64         `json:"pauseDurationMillis"`
	Stages              []JenkinsStages `json:"stages"`
}
