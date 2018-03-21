# Build pysssix driver
FROM golang:1.9.2

COPY drivers/pysssix/main.go /go
RUN go build /go/main.go


# Build goofys driver
FROM golang:1.9.2

COPY drivers/goofys/main.go /go
RUN go build /go/main.go


# Build deployment container
FROM bash:4.4

COPY deploy.sh /usr/local/bin
COPY --from=0 /go/main /pysssix-flex-volume
COPY --from=1 /go/main /goofys-flex-volume

CMD /usr/local/bin/deploy.sh
