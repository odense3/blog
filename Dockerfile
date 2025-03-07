FROM golang:1.23.6-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

FROM gcr.io/distroless/base-debian10

COPY --from=builder /app/main /app/main
COPY ./docs /app/docs

WORKDIR /app

EXPOSE 8000

CMD ["/app/main"]