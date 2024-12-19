FROM docker.io/grafana/k6:0.55.1@sha256:88a0ce2faa0f4b80e7b49191d13d86da2cb5bbf302696f93e1f68bb6c2fd0400 as k6-image

FROM alpine:3.21@sha256:21dc6063fd678b478f57c0e13f47560d0ea4eeba26dfc947b2a4f81f686b9f45 as certs

RUN apk add --update --no-cache ca-certificates bash

COPY --from=k6-image /usr/bin/k6 /usr/bin/k6

COPY release/vela-k6 /bin/vela-k6

ENTRYPOINT ["/bin/vela-k6"]
