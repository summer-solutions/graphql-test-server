FROM golang:1.15-buster as builder

WORKDIR /app

COPY . ./

RUN go build -mod=readonly -v -o accounts ./services/accounts
RUN go build -mod=readonly -v -o products ./services/products
RUN go build -mod=readonly -v -o reviews ./services/reviews

FROM debian:buster-slim
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    --no-install-recommends \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/accounts /app/accounts
COPY --from=builder /app/products /app/products
COPY --from=builder /app/reviews /app/reviews

WORKDIR /app