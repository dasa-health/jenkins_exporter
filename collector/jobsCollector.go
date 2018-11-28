package collector

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	logger "github.com/dasa-health/elk-logger"
	"github.com/prometheus/client_golang/prometheus"
)

// JobsMetricsCollector create collector metrics
func JobsMetricsCollector() *JobsCollector {
	return &JobsCollector{
		job: prometheus.NewDesc("jenkins_metrics_jobs_total",
			"Shows metrics from job",
			[]string{"count"}, nil,
		),
	}
}

// Describe implemented with dummy data to satisfy interface.
func (c *JobsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- prometheus.NewDesc("dummy", "dummy", nil, nil)
}

//Collect implements required collect function for all promehteus collectors
func (c *JobsCollector) Collect(ch chan<- prometheus.Metric) {

	logger.Info("Start exporter from get all jobs")

	jenkins := getAllJobs()

	if len(jenkins.Jobs) <= 0 {
		logger.Error("Job %s not found on Jenkins")
	}

	logger.Info("Jobs returned from Filter", jenkins.Jobs)
	createdTimestamp := float64(time.Now().UnixNano() / 1000000)

	ch <- prometheus.MustNewConstMetric(c.job,
		prometheus.CounterValue,
		createdTimestamp,
		fmt.Sprintf("%d", len(jenkins.Jobs)))

	logger.Info("Finish exporter from get all jobs")
}

func getAllJobs() Jenkins {
	client := http.Client{}
	request, err := http.NewRequest("GET", os.Getenv("jenkinsHost")+"/api/json", nil)
	request.SetBasicAuth(os.Getenv("jenkinsUser"), os.Getenv("jenkinsPassword"))

	res, err := client.Do(request)

	if err != nil {
		logger.Error(fmt.Sprintf("[getAllJobs] - Error in GET %s", request.URL), err)
		return Jenkins{}
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		logger.Error(fmt.Sprintf("[getAllJobs] - ReadError in GET %s", request.URL), readErr)
		return Jenkins{}
	}

	var result Jenkins
	jsonError := json.Unmarshal(body, &result)
	if jsonError != nil {
		logger.Error(fmt.Sprintf("[getAllJobs] - JsonError in GET %s", request.URL), jsonError)
		return Jenkins{}
	}

	return result
}

func filterJobs(n []Jobs) []Jobs {
	var result []Jobs
	for _, x := range n {
		if x.Type == "org.jenkinsci.plugins.workflow.job.WorkflowJob" {
			result = append(result, x)
		}
	}
	return result
}
