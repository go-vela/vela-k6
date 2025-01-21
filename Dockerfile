FROM docker.io/grafana/k6:0.56.0@sha256:89684628f98358ba6cc1a2e604bf2b05d49aad43611ece73a44aafcec06fcf28 as k6-image

FROM alpine:3.21@sha256:56fa17d2a7e7f168a043a2712e63aed1f8543aeafdcee47c58dcffe38ed51099 as certs

RUN apk add --update --no-cache ca-certificates bash

COPY --from=k6-image /usr/bin/k6 /usr/bin/k6

COPY release/vela-k6 /bin/vela-k6

ENTRYPOINT ["/bin/vela-k6"]
