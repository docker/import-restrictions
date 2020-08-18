# syntax=docker/dockerfile:experimental


#   Copyright 2020 Docker, Inc.

#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at

#       http://www.apache.org/licenses/LICENSE-2.0

#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.

ARG GO_VERSION=1.15.0-alpine
ARG GOLANGCI_LINT_VERSION=v1.30.0-alpine

FROM --platform=${BUILDPLATFORM} golang:${GO_VERSION} AS base
RUN apk add --no-cache make
WORKDIR /import-restrictions
ENV GO111MODULE=on
COPY go.* .
RUN go mod download

FROM base AS make-build
ARG TARGETOS
ARG TARGETARCH
ENV CGO_ENABLED=0
RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH} \
    make -f builder.Makefile build

FROM scratch AS build
COPY --from=make-build /out/* .
