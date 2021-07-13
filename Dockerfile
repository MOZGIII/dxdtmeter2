FROM golang:1.17beta1 AS builder

WORKDIR /build
COPY main.go .
RUN CGO_ENABLED=0 go build -o app main.go

FROM scratch

COPY --from=builder /build/app /app

CMD ["/app"]
