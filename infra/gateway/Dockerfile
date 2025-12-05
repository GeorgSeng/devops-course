# Stage 1: Build the Go application
FROM golang:1.25 as builder

WORKDIR /app
COPY . .
#COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -o gateway .

# Stage 2: Create the final image with NGINX and Supervisor
FROM nginx:latest

# Install supervisor
RUN apt-get update && apt-get install -y supervisor && rm -rf /var/lib/apt/lists/*

# Copy NGINX configuration
COPY nginx.conf /etc/nginx/nginx.conf

# Copy supervisord configuration
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf

# Copy the Go application from the builder stage
COPY --from=builder /app/gateway /usr/local/bin/gateway

# Expose port 80 for NGINX
EXPOSE 80
# Expose port 8080 for the Golang API
EXPOSE 8080
# Command to run supervisord
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]
