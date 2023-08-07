FROM docker.io/grafana/k6:latest as k6-image

FROM alpine as certs

RUN apk add --update --no-cache ca-certificates

COPY --from=k6-image /usr/bin/k6 /usr/bin/k6

COPY release/vela-k6 /bin/vela-k6

ENTRYPOINT ["/bin/vela-k6"]
