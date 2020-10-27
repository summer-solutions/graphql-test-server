FROM debian:buster-slim as builder
WORKDIR /app
#RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
#    ca-certificates && \
#    rm -rf /var/lib/apt/lists/*
COPY server ./
COPY --from=builder server /app/server

CMD ["/app/server"]