FROM golang:1.8.3-alpine3.6

WORKDIR /go/src/app

RUN apk --no-cache add ca-certificates curl git
RUN curl -s -L https://github.com/golang/dep/releases/download/v0.3.1/dep-linux-amd64 -o /usr/bin/dep && chmod +x /usr/bin/dep

ADD *.go Gopkg.lock Gopkg.toml ./

RUN dep ensure
RUN go-wrapper install

CMD [ "go-wrapper", "run" ]

