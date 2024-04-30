# Step 1: Modules
FROM golang:1.21.2-alpine3.17 as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:1.21.2-alpine3.17 as builder
COPY --from=modules /go/pkg /go/pkg
COPY . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o bot ./cmd/bot/main.go

# Step 3: Final
FROM alpine:3.16
COPY --from=builder /build/bot .
CMD [ "./bot", "migrate"]