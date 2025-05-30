FROM docker.io/grafana/k6:0.59.0@sha256:654b5a04672361b7f5ed76359c985e0b2a1e28e0ca15495466ad4c732006e1f3 as k6-image

FROM alpine:3.21@sha256:a8560b36e8b8210634f77d9f7f9efd7ffa463e380b75e2e74aff4511df3ef88c as certs

RUN apk add --update --no-cache ca-certificates bash

COPY --from=k6-image /usr/bin/k6 /usr/bin/k6

COPY release/vela-k6 /bin/vela-k6

ENTRYPOINT ["/bin/vela-k6"]
