FROM golang:1.24.4-alpine AS final
WORKDIR /app
COPY . .
RUN go build -o main ./cmd/server
RUN chmod +x start.sh wait-for.sh

EXPOSE 8080
CMD ["/app/main"]
