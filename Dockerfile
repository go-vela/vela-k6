# Copyright (c) 2023 Target Brands, Inc. All rights reserved.
#
# Use of this source code is governed by the LICENSE file in this repository.

FROM docker.io/grafana/k6:latest as k6-image

FROM alpine as certs

RUN apk add --update --no-cache ca-certificates

COPY --from=k6-image /usr/bin/k6 /usr/bin/k6

COPY vela-k6 /bin/vela-k6

ENTRYPOINT ["/bin/vela-k6"]
