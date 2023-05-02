FROM golang:latest AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN make build

FROM scratch
WORKDIR /app
COPY --from=builder /app/build/userapi .
COPY migrations migrations
EXPOSE 8000

CMD ["./userapi"]
