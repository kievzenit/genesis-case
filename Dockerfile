FROM golang:1.24-alpine AS builder
ARG CGO_ENABLED=0
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o server ./cmd/server/main.go

FROM scratch
COPY --from=builder /app/server /server
COPY ./migrations /migrations
COPY ./templates /templates

EXPOSE 8080

ENTRYPOINT ["/server"]