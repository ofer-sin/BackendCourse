# # build stage
# FROM golang:1.24-alpine3.21 AS builder
# WORKDIR /app
# COPY . .
# RUN go build -o main main.go
# RUN apk add curl
# RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.3/migrate.linux-amd64.tar.gz | tar xvz

# # run stage. Copy only the exceutable, no need for the whole goland packages
# # start from a clean alpine image
# FROM alpine:3.21
# WORKDIR /app
# COPY --from=builder /app/main .
# COPY --from=builder /app/migrate ./migrate
# COPY app.env .
# COPY start.sh .
# COPY wait-for.sh .
# COPY db/migration ./migration

# # espose the port the app runs on
# EXPOSE 8080

# # run the executable
# CMD [ "/app/main" ]
# ENTRYPOINT [ "/app/start.sh" ]


# build stage
FROM golang:1.24-alpine3.21 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.3/migrate.linux-amd64.tar.gz | tar xvz

# run stage. Copy only the executable, no need for the whole golang packages
# start from a clean alpine image
FROM alpine:3.21
WORKDIR /app

# Install necessary packages
RUN apk add --no-cache bash

# Copy files from builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY app.env .
COPY start.sh .
# COPY wait-for.sh .
COPY db/migration ./migration

# Make scripts executable
RUN chmod +x start.sh
# RUN chmod +x wait-for.sh

# Expose the port the app runs on
EXPOSE 8080

# Run the executable
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]
