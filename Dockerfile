FROM golang:1.20-alpine AS builder
WORKDIR /work
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go mod tidy
RUN ls
RUN go build -o /work/backend ./cmd/server

FROM alpine:3.14
COPY --from=builder /work/backend /work/backend
EXPOSE 3001
ENTRYPOINT [ "/work/backend" ]