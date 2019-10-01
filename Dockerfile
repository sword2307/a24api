FROM golang:alpine as builder

RUN apk update && apk --update add git && apk --update add upx && apk --update add ca-certificates
RUN CGO_ENABLED=0 go get -a -ldflags '-s -w' github.com/sword2307/a24api && upx --brute /go/bin/a24api

FROM scratch
COPY --from=builder /go/bin/a24api /srv/a24api
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENTRYPOINT ["/srv/a24api"]
