FROM golang:1.12-alpine as builder
ARG version
ARG gitCommit


COPY ./ /go/src/goproxy
WORKDIR /go/src/goproxy

RUN apk update && \
    apk add -U git wget && \
    GO111MODULE=on \
    CGO_ENABLED=0 go build -a -tags netgo -installsuffix netgo \
    -o /app/goproxy && \
    chmod u+x /app/goproxy


FROM golang:1.12-alpine

ARG gitCommit
ENV GIT_COMMIT=${gitCommit}
COPY --from=builder /app/ /app/

RUN apk update && \
    apk add -U git wget

WORKDIR /app/
CMD ["./goproxy"]
