FROM docker.io/grafana/k6:latest@sha256:10683d5b3803e61567476136b6dcf4ac06df67bdc4a1c9fbbeab61bc916ccd77 as k6-image

FROM alpine@sha256:7144f7bab3d4c2648d7e59409f15ec52a18006a128c733fcff20d3a4a54ba44a as certs

RUN apk add --update --no-cache ca-certificates

COPY --from=k6-image /usr/bin/k6 /usr/bin/k6

COPY release/vela-k6 /bin/vela-k6

ENTRYPOINT ["/bin/vela-k6"]
