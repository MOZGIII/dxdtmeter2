FROM golang:1.15 AS builder

WORKDIR /build
COPY main.go .
RUN go build -o app main.go

FROM scratch

COPY --from=builder /build/app /app

CMD ["/app"]
