# Jenkins exporter

The Jenkins exporter Prometheus

## Building and running

### Local Build

    export port=3000 # optional variable, its default value is 3000
    export jenkinsHost=http://your_host.com
    export jenkinsUser=your_user
    export jenkinsPassword=your_pass
    export environment=prd 
    export release=1.0.0 
    export elkHost=http://localhost:9200 
    export elkIndex=jenkins-exporter 
    export activeLogSegregation=true 

    go run main.go

See URL [http://localhost:3000/jobs](http://localhost:3000/jobs) to retrieve metrics from count all jobs from jenkins

See URL [http://localhost:3000/metrics?jobname=your_jobname](http://localhost:3000/metrics?jobname=your_jobname) to retrieve metrics from all pipeline type jobs to Jenkins code that were preconfigured.

### Building with Docker

    docker build --tag jenkins_exporter --build-arg port=3000 --build-arg jenkinsHost=http://your_host.com --build-arg jenkinsUser=your_user --build-arg jenkinsPassword=your_pass --build-arg environment=prd --build-arg release=1.0.0 --build-arg elkHost=http://localhost:9200 --build-arg elkIndex=jenkins-exporter --build-arg activeLogSegregation=true -f Dockerfile .

    docker run -d -p 3000:3000 --name jenkins_exporter 

## Releases

### 2.0.0 (28.11.2018)

* Standardization of environment variables for camel case
* Inclusion of log segregation for the elasticsearch
* Creating the route that returns total jobs in jenkins
* Changing the main functionality of the exporter, where instead of searching all the jobs and the metrics of each 1, we search the data of a specific job via query string

### 1.1.0 (08.10.2018)

* Automatically Recovering All Jenkins Type Pipeline as Code Jobs
* Removing environment variables JENKINS_JOBS_SEPARATOR and JENKINS_JOBS
* Standardizing the name and value of fields from milliseconds to seconds

### 1.0.0 (05.10.2018)

* Prometheus exporter's initial version that retrieves the statuses of jenkins's pipeline type jobs