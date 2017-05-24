FROM golang:1.7
ADD . /go/src/ipupdate
WORKDIR /go/src/ipupdate
RUN cd /go/src/ipupdate && go get ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ipupdate .

FROM scratch
WORKDIR /root/
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=0 /go/src/ipupdate/ipupdate .
CMD ["./ipupdate"]
