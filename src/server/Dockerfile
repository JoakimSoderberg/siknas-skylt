FROM golang:1.8.3-alpine3.6

WORKDIR /go/src/app

# We install Java/xvfb to enable running Processing sketches in headless mode.
RUN apk --no-cache add ca-certificates curl git \
    openjdk8-jre xvfb libxrender libxtst libxi

RUN curl -s -L https://github.com/golang/dep/releases/download/v0.3.1/dep-linux-amd64 -o /usr/bin/dep && chmod +x /usr/bin/dep

ADD *.go Gopkg.lock Gopkg.toml ./
RUN dep ensure -v
RUN go-wrapper install

ENV DISPLAY=:1

EXPOSE 8080
EXPOSE 3000
