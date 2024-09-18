# Stage 1: Build the Go application
FROM golang:1.22.6-alpine AS builder

# Install build tools and SQLite3 dependencies
RUN apk add --no-cache gcc musl-dev go libc-dev sqlite-dev

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to cache dependencies
COPY go.mod ./

# Download dependencies
RUN go mod tidy

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=1 GOOS=linux go build -v -o url-shortener .

# Stage 2: Create a lightweight image with only the binary
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/url-shortener .

# Set build-time arguments (ARG) with defaults
ARG PORT=8080
ARG SUPERSCRT=$uperTok3nCre4te6
ARG DOMAIN=localhost

# Expose these ARG values as ENV variables for runtime
ENV PORT=${PORT}
ENV SUPERSCRT=${SUPERSCRT}
ENV DOMAIN=${DOMAIN}
ENV HTTPS=${HTTPS}
ENV LOCAL=${LOCAL}

# Expose the application port
EXPOSE ${PORT}

# Command to run the binary
CMD ["./url-shortener", "-port", "${PORT}", "-token", "${SUPERSCRT}", "-domain", "${DOMAIN}", "-https", "${HTTPS}", "-local", "${LOCAL}"]
