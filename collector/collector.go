package collector

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

// JenkinsMetricsCollector format metrics for prometheus
func JenkinsMetricsCollector() *JenkinsMetrics {

	return &JenkinsMetrics{

		job: prometheus.NewDesc("jenkins_jobs",
			"Shows metrics from job",
			[]string{"id", "name", "status", "startTime", "endTime", "durationMillis", "queueDurationMillis", "pauseDurationMillis"}, nil,
		),
		stage: prometheus.NewDesc("jenkins_jobs_stages",
			"Shows metrics from job",
			[]string{"job_id", "job_name", "id", "name", "status", "startTime", "durationMillis", "pauseDurationMillis", "execNode"}, nil,
		),
	}
}

// Describe writes all descriptors to the prometheus desc channel.
func (collector *JenkinsMetrics) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.job
	ch <- collector.stage
}

//Collect implements required collect function for all promehteus collectors
func (collector *JenkinsMetrics) Collect(ch chan<- prometheus.Metric) {

	jobs := jenkinsRequest("sre-backend-java")

	for _, job := range jobs {
		createdTimestamp := job.StartTimeMillis
		ch <- prometheus.MustNewConstMetric(collector.job,
			prometheus.CounterValue,
			createdTimestamp,
			job.ID,
			job.Name,
			job.Status,
			convert.NanoTimestampToString(job.StartTimeMillis),
			NanoTimestampToString(job.EndTimeMillis),
			strconv.FormatFloat(MillisToSecond(job.DurationMillis), 'f', -1, 64),
			strconv.FormatFloat(MillisToSecond(job.QueueDurationMillis), 'f', -1, 64),
			strconv.FormatFloat(MillisToSecond(job.PauseDurationMillis), 'f', -1, 64))

		for _, stage := range job.Stages {

			createdTimestamp := stage.StartTimeMillis
			ch <- prometheus.MustNewConstMetric(collector.stage,
				prometheus.CounterValue,
				createdTimestamp,
				job.ID,
				"sre-backend-java",
				stage.ID,
				stage.Name,
				stage.Status,
				NanoTimestampToString(stage.StartTimeMillis),
				strconv.FormatFloat(MillisToSecond(stage.DurationMillis), 'f', -1, 64),
				strconv.FormatFloat(MillisToSecond(stage.PauseDurationMillis), 'f', -1, 64),
				stage.ExecNode)
		}
	}
}

func jenkinsRequest(job string) []JenkinsJob {

	client := http.Client{}
	request, err := http.NewRequest("GET", os.Getenv("JENKINS_HOST")+"/job/"+job+"/wfapi/runs", nil)
	request.SetBasicAuth(os.Getenv("JENKINS_USER"), os.Getenv("JENKINS_PASSWORD"))

	res, err := client.Do(request)

	if err != nil {
		log.Fatalln(err)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatalln(readErr)
	}

	var result []JenkinsJob
	jsonError := json.Unmarshal(body, &result)
	if jsonError != nil {
		log.Fatalln(jsonError)
	}

	return result
}
