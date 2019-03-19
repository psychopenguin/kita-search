FROM golang:1.12 as builder

WORKDIR /usr/src/kita-search
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o kita-search .

FROM busybox:latest  
COPY --from=builder /usr/src/kita-search/kita-search .
ENTRYPOINT ["./kita-search"]
