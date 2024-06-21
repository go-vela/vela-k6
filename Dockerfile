FROM docker.io/grafana/k6:0.51.0@sha256:b2de44f593444feca8c66e7ee8cd230b272f614a7d7b2e784e7b76569c00fbf7 as k6-image

FROM alpine:3.20@sha256:b89d9c93e9ed3597455c90a0b88a8bbb5cb7188438f70953fede212a0c4394e0 as certs

RUN apk add --update --no-cache ca-certificates bash

COPY --from=k6-image /usr/bin/k6 /usr/bin/k6

COPY release/vela-k6 /bin/vela-k6

ENTRYPOINT ["/bin/vela-k6"]
