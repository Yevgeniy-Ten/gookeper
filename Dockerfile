FROM golang:1.24-alpine

WORKDIR /app


RUN apk add --no-cache gcc musl-dev


COPY go.mod ./
COPY go.sum ./


RUN go mod download


COPY . .


RUN go build -o main ./cmd/server


EXPOSE 8080

# Run the application
CMD ["./main"] 