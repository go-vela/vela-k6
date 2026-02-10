FROM docker.io/grafana/k6:1.6.0@sha256:be2a9ceb1b6ffc573277af4157135882a6d6433968b41e858ab02c5fd5847532 as k6-image

FROM alpine:3.23@sha256:25109184c71bdad752c8312a8623239686a9a2071e8825f20acb8f2198c3f659 as certs

RUN apk add --update --no-cache ca-certificates bash

COPY --from=k6-image /usr/bin/k6 /usr/bin/k6

COPY release/vela-k6 /bin/vela-k6

ENTRYPOINT ["/bin/vela-k6"]
