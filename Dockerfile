FROM docker.io/grafana/k6:1.2.3@sha256:4f82892217f3110cb233e2b2622bcc97fabc70f14bd241fbfbfe7305105c68aa as k6-image

FROM alpine:3.22@sha256:4bcff63911fcb4448bd4fdacec207030997caf25e9bea4045fa6c8c44de311d1 as certs

RUN apk add --update --no-cache ca-certificates bash

COPY --from=k6-image /usr/bin/k6 /usr/bin/k6

COPY release/vela-k6 /bin/vela-k6

ENTRYPOINT ["/bin/vela-k6"]
