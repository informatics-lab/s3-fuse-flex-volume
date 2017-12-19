FROM golang:1.9.2

COPY main.go /go
RUN go build /go/main.go

FROM alpine:3.6

COPY deploy.sh /usr/local/bin
COPY --from=0 /go/main /s3-fuse-flex-volume

CMD /usr/local/bin/deploy.sh
