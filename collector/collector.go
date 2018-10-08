package collector

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/silva-willian/jenkins_exporter/services"
)

// JenkinsMetricsCollector format metrics for prometheus
func JenkinsMetricsCollector() *JenkinsMetrics {

	return &JenkinsMetrics{

		job: prometheus.NewDesc("jenkins_jobs",
			"Shows metrics from job",
			[]string{"id", "name", "status", "startTime", "endTime", "durationSeconds", "queueDurationSeconds", "pauseDurationSeconds"}, nil,
		),
		stage: prometheus.NewDesc("jenkins_jobs_stages",
			"Shows metrics from stages job",
			[]string{"job_id", "job_name", "id", "name", "status", "startTime", "durationSeconds", "pauseDurationSeconds", "execNode"}, nil,
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

	jenkins := getAllJobs()

	if len(jenkins.Jobs) <= 0 {
		return
	}

	jobsSanitize := sanitizeJobs(jenkins.Jobs)

	if len(jobsSanitize) <= 0 {
		return
	}

	for _, jenkinsJob := range jobsSanitize {

		jobs := getJobByName(jenkinsJob.Name)

		for _, job := range jobs {
			createdTimestamp := job.StartTimeMillis
			ch <- prometheus.MustNewConstMetric(collector.job,
				prometheus.CounterValue,
				createdTimestamp,
				job.ID,
				jenkinsJob.Name,
				job.Status,
				services.NanoTimestampToString(job.StartTimeMillis),
				services.NanoTimestampToString(job.EndTimeMillis),
				strconv.FormatFloat(services.MillisToSecond(job.DurationMillis), 'f', -1, 64),
				strconv.FormatFloat(services.MillisToSecond(job.QueueDurationMillis), 'f', -1, 64),
				strconv.FormatFloat(services.MillisToSecond(job.PauseDurationMillis), 'f', -1, 64))

			for _, stage := range job.Stages {

				createdTimestamp := stage.StartTimeMillis
				ch <- prometheus.MustNewConstMetric(collector.stage,
					prometheus.CounterValue,
					createdTimestamp,
					job.ID,
					jenkinsJob.Name,
					stage.ID,
					stage.Name,
					stage.Status,
					services.NanoTimestampToString(stage.StartTimeMillis),
					strconv.FormatFloat(services.MillisToSecond(stage.DurationMillis), 'f', -1, 64),
					strconv.FormatFloat(services.MillisToSecond(stage.PauseDurationMillis), 'f', -1, 64),
					stage.ExecNode)
			}
		}
	}
}

func getJobByName(job string) []JobDetails {

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

	var result []JobDetails
	jsonError := json.Unmarshal(body, &result)
	if jsonError != nil {
		log.Fatalln(jsonError)
	}

	return result
}

func getAllJobs() Jenkins {
	client := http.Client{}
	request, err := http.NewRequest("GET", os.Getenv("JENKINS_HOST")+"/api/json", nil)
	request.SetBasicAuth(os.Getenv("JENKINS_USER"), os.Getenv("JENKINS_PASSWORD"))

	res, err := client.Do(request)

	if err != nil {
		log.Fatalln(err)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatalln(readErr)
	}

	var result Jenkins
	jsonError := json.Unmarshal(body, &result)
	if jsonError != nil {
		log.Fatalln(jsonError)
	}

	return result
}

func sanitizeJobs(n []Jobs) []Jobs {
	var result []Jobs
	for _, x := range n {
		if x.Type == "org.jenkinsci.plugins.workflow.job.WorkflowJob" {
			result = append(result, x)
		}
	}
	return result
}
