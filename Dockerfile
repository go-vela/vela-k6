FROM docker.io/grafana/k6:0.52.0@sha256:895160792a74382ba8c0be3ca51a3acff8becf19cd566bbd1d8d43b3cbfa3a73 as k6-image

FROM alpine:3.20@sha256:0a4eaa0eecf5f8c050e5bba433f58c052be7587ee8af3e8b3910ef9ab5fbe9f5 as certs

RUN apk add --update --no-cache ca-certificates bash

COPY --from=k6-image /usr/bin/k6 /usr/bin/k6

COPY release/vela-k6 /bin/vela-k6

ENTRYPOINT ["/bin/vela-k6"]
