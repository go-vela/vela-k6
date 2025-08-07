FROM docker.io/grafana/k6:1.1.0@sha256:aa8202f377550cee0c8bad295bbe8d2d4d4cf88d15c98383e9ecc53c56882308 as k6-image

FROM alpine:3.22@sha256:8a1f59ffb675680d47db6337b49d22281a139e9d709335b492be023728e11715 as certs

RUN apk add --update --no-cache ca-certificates bash

COPY --from=k6-image /usr/bin/k6 /usr/bin/k6

COPY release/vela-k6 /bin/vela-k6

ENTRYPOINT ["/bin/vela-k6"]
