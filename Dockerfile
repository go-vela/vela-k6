FROM docker.io/grafana/k6:0.47.0@sha256:f7650fff23d183b96a51be366aaadd90a29548585a24407f4ec9712cbd66ec73 as k6-image

FROM alpine:3.18@sha256:eece025e432126ce23f223450a0326fbebde39cdf496a85d8c016293fc851978 as certs

RUN apk add --update --no-cache ca-certificates bash

COPY --from=k6-image /usr/bin/k6 /usr/bin/k6

COPY release/vela-k6 /bin/vela-k6

ENTRYPOINT ["/bin/vela-k6"]
