FROM golang:1.9.2

COPY main.go /go
RUN go build /go/main.go

FROM alpine:3.6

COPY deploy.sh /usr/local/bin
COPY --from=0 /go/main /nfs-flex-volume
COPY s3fuse.py /usr/local/src/s3fuse.py
COPY s3fuseenv.pex /usr/local/src/s3fuseenv.pex

CMD /usr/local/bin/deploy.sh
