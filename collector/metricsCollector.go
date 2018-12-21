package collector

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	logger "github.com/dasa-health/elk-logger"
	"github.com/dasa-health/jenkins_exporter/services"
	"github.com/prometheus/client_golang/prometheus"
)

// JenkinsMetricsCollector create collector metrics
func JenkinsMetricsCollector(name string) *JenkinsCollector {
	return &JenkinsCollector{
		name: name,
		job: prometheus.NewDesc("jenkins_metrics_job_total",
			"Shows metrics from job",
			[]string{"id", "name", "status", "startTime", "endTime", "durationMillis", "queueDurationMillis", "pauseDurationMillis"}, nil,
		),
		stage: prometheus.NewDesc("jenkins_metrics_job_stages_total",
			"Shows metrics from stages job",
			[]string{"job_id", "job_name", "id", "name", "status", "startTime", "durationMillis", "pauseDurationMillis", "execNode"}, nil,
		),
	}
}

// Describe implemented with dummy data to satisfy interface.
func (c *JenkinsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- prometheus.NewDesc("dummy", "dummy", nil, nil)
}

//Collect implements required collect function for all promehteus collectors
func (c *JenkinsCollector) Collect(ch chan<- prometheus.Metric) {

	logger.Info(fmt.Sprintf("Start exporter from job name %s", c.name))

	if c.name == "" {
		logger.Error("Job name is empty")
		return
	}

	jobs := getJobByName(c.name)

	if jobs == nil || len(jobs) <= 0 {
		logger.Error(fmt.Sprintf("Job %s not found on Jenkins", c.name))
	}

	logger.Info("Jobs returned from Jenkins", jobs)

	for _, job := range jobs {
		createdTimestamp := job.StartTimeMillis
		ch <- prometheus.MustNewConstMetric(c.job,
			prometheus.CounterValue,
			createdTimestamp,
			job.ID,
			c.name,
			job.Status,
			services.NanoTimestampToString(job.StartTimeMillis),
			services.NanoTimestampToString(job.EndTimeMillis),
			strconv.FormatFloat(job.DurationMillis, 'f', -1, 64),
			strconv.FormatFloat(job.QueueDurationMillis, 'f', -1, 64),
			strconv.FormatFloat(job.PauseDurationMillis, 'f', -1, 64))

		for _, stage := range job.Stages {
			createdTimestamp := stage.StartTimeMillis
			ch <- prometheus.MustNewConstMetric(c.stage,
				prometheus.CounterValue,
				createdTimestamp,
				job.ID,
				c.name,
				stage.ID,
				stage.Name,
				stage.Status,
				services.NanoTimestampToString(stage.StartTimeMillis),
				strconv.FormatFloat(stage.DurationMillis, 'f', -1, 64),
				strconv.FormatFloat(stage.PauseDurationMillis, 'f', -1, 64),
				stage.ExecNode)
		}
	}

	logger.Info(fmt.Sprintf("Finish exporter from job name %s", c.name))
}

func getJobByName(job string) []JobDetails {

	client := http.Client{}
	request, err := http.NewRequest("GET", os.Getenv("jenkinsHost")+"/job/"+job+"/wfapi/runs", nil)
	request.SetBasicAuth(os.Getenv("jenkinsUser"), os.Getenv("jenkinsPassword"))

	res, err := client.Do(request)

	if err != nil {
		logger.Error(fmt.Sprintf("[getJobByName] - Error in GET %s", request.URL), err)
		return nil
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		logger.Error(fmt.Sprintf("[getJobByName] - ReadError in GET %s", request.URL), readErr)
		return nil
	}

	var result []JobDetails
	jsonError := json.Unmarshal(body, &result)
	if jsonError != nil {
		logger.Error(fmt.Sprintf("[getJobByName] - JsonError in GET %s", request.URL), jsonError)
		return nil
	}

	return result
}
