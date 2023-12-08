FROM docker.io/grafana/k6:0.48.0@sha256:8f3404658de9c66ba4c446c21bf88beeac679fdaacee9245d21d68a5c7930d39 as k6-image

FROM alpine:3.19@sha256:51b67269f354137895d43f3b3d810bfacd3945438e94dc5ac55fdac340352f48 as certs

RUN apk add --update --no-cache ca-certificates bash

COPY --from=k6-image /usr/bin/k6 /usr/bin/k6

COPY release/vela-k6 /bin/vela-k6

ENTRYPOINT ["/bin/vela-k6"]
