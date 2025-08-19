FROM docker.io/grafana/k6:1.2.2@sha256:a98db15eb83dfc4dbded9653a15f53557c011400019eba715a9cd15b0ab709d9 as k6-image

FROM alpine:3.22@sha256:4bcff63911fcb4448bd4fdacec207030997caf25e9bea4045fa6c8c44de311d1 as certs

RUN apk add --update --no-cache ca-certificates bash

COPY --from=k6-image /usr/bin/k6 /usr/bin/k6

COPY release/vela-k6 /bin/vela-k6

ENTRYPOINT ["/bin/vela-k6"]
