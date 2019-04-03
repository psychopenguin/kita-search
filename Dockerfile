FROM golang:1.12 AS builder
LABEL name="kita-search" \
    version=0.0.1
RUN apt-get update && apt-get install -y upx && apt-get clean
WORKDIR /usr/src/kita-search
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -v -o kita-search .
RUN upx --brute kita-search

FROM busybox:1
COPY --from=builder /usr/src/kita-search/kita-search .
CMD ["./kita-search"]
EXPOSE 5000
