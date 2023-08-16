FROM docker.io/grafana/k6:0.46.0@sha256:2f40a302ec1e1e3cc96b9a3871bf5d7d4697e9ecc4fe90546ba0eb005d3458e3 as k6-image

FROM alpine:3.18@sha256:7144f7bab3d4c2648d7e59409f15ec52a18006a128c733fcff20d3a4a54ba44a as certs

RUN apk add --update --no-cache ca-certificates bash

COPY --from=k6-image /usr/bin/k6 /usr/bin/k6

COPY release/vela-k6 /bin/vela-k6

ENTRYPOINT ["/bin/vela-k6"]
