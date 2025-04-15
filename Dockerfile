FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY ./ ./

RUN go build ./cmd/backend/main.go

FROM alpine:latest AS runner

COPY --from=build /app/main /

CMD ["/main"]

