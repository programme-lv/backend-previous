FROM golang:1.19 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -v -o main ./cmd/server

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 3001
CMD ["./main"] 
