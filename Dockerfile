FROM golang:1.10-alpine

LABEL maintainer="Dasa <devops-team@dasa.com.br>"

ARG port
ARG jenkinsHost
ARG jenkinsUser
ARG jenkinsPassword
ARG environment
ARG release
ARG elkHost
ARG elkIndex
ARG activeLogSegregation

ENV port $port
ENV jenkinsHost $jenkinsHost
ENV jenkinsUser $jenkinsUser
ENV jenkinsPassword $jenkinsPassword
ENV environment $environment
ENV release $release
ENV elkHost $elkHost
ENV elkIndex $elkIndex
ENV activeLogSegregation $activeLogSegregation

WORKDIR /go/src/app
COPY . .

RUN apk add util-linux && apk add --update git

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE $PORT
CMD ["app"]