FROM docker.io/grafana/k6:0.53.0@sha256:9b48ed1865a697ac481b5aed98975d098c6086ad418a023eed12e43d876ce271 as k6-image

FROM alpine:3.20@sha256:beefdbd8a1da6d2915566fde36db9db0b524eb737fc57cd1367effd16dc0d06d as certs

RUN apk add --update --no-cache ca-certificates bash

COPY --from=k6-image /usr/bin/k6 /usr/bin/k6

COPY release/vela-k6 /bin/vela-k6

ENTRYPOINT ["/bin/vela-k6"]
