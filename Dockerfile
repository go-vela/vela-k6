FROM docker.io/grafana/k6:1.5.0@sha256:2072ea9eafa596532d9aee0cc0e0a50cfb0e7fb703981a46179af5f4f22c5ea4 as k6-image

FROM alpine:3.23@sha256:865b95f46d98cf867a156fe4a135ad3fe50d2056aa3f25ed31662dff6da4eb62 as certs

RUN apk add --update --no-cache ca-certificates bash

COPY --from=k6-image /usr/bin/k6 /usr/bin/k6

COPY release/vela-k6 /bin/vela-k6

ENTRYPOINT ["/bin/vela-k6"]
