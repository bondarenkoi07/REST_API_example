FROM golang:1.15.7-buster

RUN go get -u github.com/bondarenkoi07/REST_API_example

WORKDIR /go/src/github.com/bondarenkoi07/REST_API_example

# RUN go mod init && \
 #   go mod download && \
  #  go mod vendor && \
 #   go mod verify && \
  #  go install

#CMD ["go","test","/go/src/github.com/bondarenkoi07/REST_API_example/database"]
CMD ["/bin/bash"]