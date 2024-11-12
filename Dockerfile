FROM docker.io/grafana/k6:0.55.0@sha256:b24f418fc99a26dd57904c952c03bfaf79462be18508acc45aafa07ff68e7df2 as k6-image

FROM alpine:3.20@sha256:1e42bbe2508154c9126d48c2b8a75420c3544343bf86fd041fb7527e017a4b4a as certs

RUN apk add --update --no-cache ca-certificates bash

COPY --from=k6-image /usr/bin/k6 /usr/bin/k6

COPY release/vela-k6 /bin/vela-k6

ENTRYPOINT ["/bin/vela-k6"]
