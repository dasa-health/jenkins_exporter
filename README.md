# Jenkins exporter

The Jenkins exporter Prometheus

## Building and running

### Local Build

    export PORT=3000 # optional variable, its default value is 3000
    export JENKINS_HOST=http://your_host.com
    export JENKINS_USER=your_user
    export JENKINS_PASSWORD=your_pass
    export JENKINS_JOBS_SEPARATOR=, # variable, its default value is ","
    export JENKINS_JOBS"=job1,job2,job3 # name of jobs to be retrieved jenkins metrics

    go run main.go

See URL [http://localhost:3000/metrics](http://localhost:3000/metrics) to retrieve metrics from all pipeline type jobs to Jenkins code that were preconfigured.

### Building with Docker

    docker build --tag jenkins_exporter --build-arg PORT=3000 --build-arg JENKINS_HOST=http://your_host.com --build-arg JENKINS_USER=your_user --build-arg JENKINS_PASSWORD=your_pass --build-arg JENKINS_JOBS_SEPARATOR=, --build-arg JENKINS_JOBS=job1,job2,job3 -f Dockerfile .

    docker run -d -p 3000:3000 --name jenkins_exporter 

## Releases

### 1.0.0 (05.10.2018)

* Prometheus exporter's initial version that retrieves the statuses of jenkins's pipeline type jobs