FROM docker.io/grafana/k6:0.47.0@sha256:f7650fff23d183b96a51be366aaadd90a29548585a24407f4ec9712cbd66ec73 as k6-image

FROM alpine:3.18@sha256:34871e7290500828b39e22294660bee86d966bc0017544e848dd9a255cdf59e0 as certs

RUN apk add --update --no-cache ca-certificates bash

COPY --from=k6-image /usr/bin/k6 /usr/bin/k6

COPY release/vela-k6 /bin/vela-k6

ENTRYPOINT ["/bin/vela-k6"]
