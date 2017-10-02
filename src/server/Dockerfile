FROM golang:1.9

WORKDIR /go/src/app

ADD https://github.com/golang/dep/releases/download/v0.3.1/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep

COPY . .

RUN dep ensure
RUN go-wrapper install    # "go install -v ./..."

EXPOSE 8080

CMD ["go-wrapper", "run"] # ["app"]

