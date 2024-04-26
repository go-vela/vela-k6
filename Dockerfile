FROM docker.io/grafana/k6:0.50.0@sha256:0a1289901ecf46819c50a7dd00e1be3e85d2b26dea175dc72f3cd730317cd584 as k6-image

FROM alpine:3.19@sha256:c5b1261d6d3e43071626931fc004f70149baeba2c8ec672bd4f27761f8e1ad6b as certs

RUN apk add --update --no-cache ca-certificates bash

COPY --from=k6-image /usr/bin/k6 /usr/bin/k6

COPY release/vela-k6 /bin/vela-k6

ENTRYPOINT ["/bin/vela-k6"]
