FROM golang:1.24 as builder

WORKDIR /app
COPY . .
RUN ls
RUN CGO_ENABLED=0 GOOS=linux go build -o myapp ./

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/myapp /app/

EXPOSE 8080
CMD ["/app/myapp"]