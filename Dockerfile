FROM golang:1.26-alpine AS builder

ARG BUILD_SHA=dev
RUN echo "BUILD_SHA during docker build: $BUILD_SHA"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make generate
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.BuildSHA=$BUILD_SHA" -o app

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/app ./app

EXPOSE 8080

CMD ["./app"]