FROM docker.io/grafana/k6:1.4.1@sha256:200d24a0770ad12761569993c723fd7d48b29fc7983ff5f976bf8b8dba4c7d21 as k6-image

FROM alpine:3.22@sha256:4b7ce07002c69e8f3d704a9c5d6fd3053be500b7f1c69fc0d80990c2ad8dd412 as certs

RUN apk add --update --no-cache ca-certificates bash

COPY --from=k6-image /usr/bin/k6 /usr/bin/k6

COPY release/vela-k6 /bin/vela-k6

ENTRYPOINT ["/bin/vela-k6"]
