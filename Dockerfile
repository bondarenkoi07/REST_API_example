FROM golang:1.15.7-buster

RUN go get  github.com/bondarenkoi07/REST_API_example

WORKDIR /go/src/github.com/bondarenkoi07/REST_API_example


RUN go mod download && \
    go mod vendor && \
    go mod verify

CMD ["/bin/bash"]