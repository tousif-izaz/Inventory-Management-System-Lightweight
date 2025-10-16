FROM golang:1.23.0 AS builder

WORKDIR /build
COPY go.* .
RUN go mod download
COPY . .
RUN go build -o ./ims cmd/imsapi/main.go

FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /build/ims ./ims
COPY --from=builder /build/schema_sqlite.sql ./schema_sqlite.sql

ENV DB_PATH=/app/data/inventory.db
ENV JWT_KEY=ultra-super-secret

CMD ["/app/ims"]
