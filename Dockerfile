# Build stage
FROM golang:1.20.3-alpine3.17 AS build

WORKDIR /build
COPY . .
RUN go build -o app .

# Final stage
FROM alpine:3.17

WORKDIR /app
COPY --from=build /build/app .
CMD ["./app"]
