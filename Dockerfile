FROM golang:1.16-alpine AS builder

WORKDIR /app

COPY . .

RUN apk add --no-cache git make
RUN rm go.sum && go mod tidy
RUN make build

FROM alpine:latest

WORKDIR /root

RUN apk --no-cache add ca-certificates

COPY --from=builder /app/build/lyre-be ./
COPY --from=builder /app/configs/dev/local_config.yaml ./
RUN mkdir tmp/

ENTRYPOINT ./lyre-be run --config=local_config.yaml
