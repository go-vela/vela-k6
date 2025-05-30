FROM docker.io/grafana/k6:1.0.0@sha256:f21270290d702cbf0a7d6ba5d7ed100b63bcb233b558b885ed787547b3488279 as k6-image

FROM alpine:3.21@sha256:a8560b36e8b8210634f77d9f7f9efd7ffa463e380b75e2e74aff4511df3ef88c as certs

RUN apk add --update --no-cache ca-certificates bash

COPY --from=k6-image /usr/bin/k6 /usr/bin/k6

COPY release/vela-k6 /bin/vela-k6

ENTRYPOINT ["/bin/vela-k6"]
