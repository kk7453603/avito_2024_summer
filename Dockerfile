FROM golang:1.24-alpine AS builder
WORKDIR /usr/local/src
RUN apk --no-cache add bash gcc musl-dev

# Dependencies
COPY ["go.mod", "go.sum", "./"]
RUN go mod download

# Build
COPY . .
RUN go build -o ./bin/app ./cmd/merchshop/main.go


FROM alpine AS runner
WORKDIR /usr/local/src

# Dependencies
COPY ./.env ./.env
COPY --from=builder /usr/local/src/bin/app .

CMD ["./app"]
