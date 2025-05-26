# build stage
FROM golang:1.24-alpine3.21 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# run stage. Copy only the exceutable, no need for the whole goland packages
# start from a clean alpine image
FROM alpine:3.21
WORKDIR /app
COPY --from=builder /app/main .

# espose the port the app runs on
EXPOSE 8080

# run the executable
CMD ["./app/main"]
