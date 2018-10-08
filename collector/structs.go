package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

// JenkinsMetrics is the struct that maps the entities to the Prometheus client
type JenkinsMetrics struct {
	job   *prometheus.Desc
	stage *prometheus.Desc
}

// JobStagesDetails is the struct that does the mapping with the stages of the Jenkis
type JobStagesDetails struct {
	Name                string  `json:"name"`
	ID                  string  `json:"id"`
	Status              string  `json:"status"`
	StartTimeMillis     float64 `json:"startTimeMillis"`
	DurationMillis      float64 `json:"durationMillis"`
	PauseDurationMillis float64 `json:"pauseDurationMillis"`
	ExecNode            string  `json:"execNode"`
}

// JobDetails is the struct that represents the jobs of Jenkins
type JobDetails struct {
	Name                string             `json:"name"`
	ID                  string             `json:"id"`
	Status              string             `json:"status"`
	StartTimeMillis     float64            `json:"startTimeMillis"`
	EndTimeMillis       float64            `json:"endTimeMillis"`
	DurationMillis      float64            `json:"durationMillis"`
	QueueDurationMillis float64            `json:"queueDurationMillis"`
	PauseDurationMillis float64            `json:"pauseDurationMillis"`
	Stages              []JobStagesDetails `json:"stages"`
}

// Jenkins is the struct that represents the Jenkins
type Jenkins struct {
	Jobs []Jobs `json:"jobs"`
}

// Jobs represents the main structure of jobs in the main Jenkis api
type Jobs struct {
	Name string `json:"name"`
	Type string `json:"_class"`
}
