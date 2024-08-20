ARG LINUX_ALPINE=alpine:3.15

FROM golang:alpine AS builder

WORKDIR /build/

COPY . .

RUN go build -o server .

FROM ${LINUX_ALPINE}

WORKDIR /app/

COPY --from=builder /build/server .
COPY --from=builder /build/config.yaml .

ENTRYPOINT ["/app/server"]