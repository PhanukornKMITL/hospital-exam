# ---------- build stage ----------
FROM golang:1.25.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Generate Swagger docs
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN /go/bin/swag init -g cmd/server/main.go

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app cmd/server/main.go


# ---------- run stage ----------
FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache tzdata
ENV TZ=Asia/Bangkok

COPY --from=builder /app/app .

EXPOSE 8080

CMD ["./app"]
