FROM docker.io/grafana/k6:0.45.0 as k6-image

FROM alpine:3.18 as certs

RUN apk add --update --no-cache ca-certificates

COPY --from=k6-image /usr/bin/k6 /usr/bin/k6

COPY release/vela-k6 /bin/vela-k6

ENTRYPOINT ["/bin/vela-k6"]
