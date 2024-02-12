FROM golang:1.21.6-alpine3.19 AS builder
LABEL authors="cukhoaimon"
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main .
COPY .env .
COPY start.sh .
COPY db/migration ./db/migration

EXPOSE 8080
ENTRYPOINT [ "/app/start.sh" ]