FROM golang:1.10-alpine

LABEL maintainer="The Prometheus Authors <prometheus-developers@googlegroups.com>"

ARG PORT
ARG JENKINS_HOST
ARG JENKINS_USER
ARG JENKINS_PASSWORD

ENV PORT $PORT
ENV JENKINS_HOST $JENKINS_HOST
ENV JENKINS_USER $JENKINS_USER
ENV JENKINS_PASSWORD $JENKINS_PASSWORD

WORKDIR /go/src/app
COPY . .

RUN apk add --update git

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE $PORT
CMD ["app"]