FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o hello

FROM scratch

COPY --from=builder /app/hello /app/hello

EXPOSE 8080

ENTRYPOINT ["/app/hello"]
