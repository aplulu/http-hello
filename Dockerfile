FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /app/hello .

FROM scratch

COPY --from=builder /app/hello /app/hello

EXPOSE 8080

ENTRYPOINT ["/app/hello"]
