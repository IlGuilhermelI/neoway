FROM golang:1.16

WORKDIR /go/src/neoway
COPY ../ /go/src/neoway/
RUN go get -d -v ./...
RUN go install -v ./...
RUN go build .
CMD ["TestNeoWay"]