FROM golang:1.26-alpine AS builder

WORKDIR /src

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/todo-api .

FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /bin/todo-api /app/todo-api
COPY conf/app.conf /app/conf/app.conf

EXPOSE 8080

CMD ["/app/todo-api"]
