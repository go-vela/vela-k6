FROM docker.io/grafana/k6:0.52.0@sha256:895160792a74382ba8c0be3ca51a3acff8becf19cd566bbd1d8d43b3cbfa3a73 as k6-image

FROM alpine:3.20@sha256:b89d9c93e9ed3597455c90a0b88a8bbb5cb7188438f70953fede212a0c4394e0 as certs

RUN apk add --update --no-cache ca-certificates bash

COPY --from=k6-image /usr/bin/k6 /usr/bin/k6

COPY release/vela-k6 /bin/vela-k6

ENTRYPOINT ["/bin/vela-k6"]
