# syntax=docker/dockerfile:1

FROM golang:1.17-alpine AS build
WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN go build -o /kpostprocess

FROM alpine:3.15
WORKDIR /
COPY --from=build /kpostprocess /kpostprocess

COPY plugin.yaml /home/argocd/cmp-server/config/plugin.yaml

RUN wget -qO- https://get.helm.sh/helm-v3.8.1-linux-amd64.tar.gz | tar zx --strip-components 1 linux-amd64/helm
RUN wget -qO- https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize%2Fv4.5.2/kustomize_v4.5.2_linux_amd64.tar.gz | tar zxv


