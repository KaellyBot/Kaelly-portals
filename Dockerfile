# Build stage
FROM golang:1.22.2-alpine3.19 AS build

WORKDIR /build
COPY . .
RUN go build -o app .

# Final stage
FROM gcr.io/distroless/static-debian12:latest

WORKDIR /app
COPY --from=build /build/app .
CMD ["./app"]
