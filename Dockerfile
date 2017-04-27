FROM golang:1.8.1-alpine

RUN mkdir -p /go/src/github.com/pjvds/slackme
WORKDIR /go/src/github.com/pjvds/slackme

COPY . /go/src/github.com/pjvds/slackme
RUN go-wrapper install \
 && rm -rf /usr/local/go \
 && rm -rf $GOPATH/pkg \
 && rm -rf $GOPATH/src

ENTRYPOINT ["cli"]
