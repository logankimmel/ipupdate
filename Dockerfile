FROM golang:1.7

ADD . /go/src/ipupdate

RUN cd /go/src/ipupdate && go get ./...
RUN go install ipupdate
WORKDIR /go/src/ipupdate
ENTRYPOINT /go/bin/ipupdate
