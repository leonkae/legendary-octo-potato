FROM golang:1.23-alpine

WORKDIR /app

# Download modules first
COPY go.mod go.sum ./
RUN go mod download

# Then copy everything else, including static/
COPY . .

RUN CGO_ENABLED=0 

# Build the application
RUN go build -o todo .

COPY  --from=builder /app/new_go .

WORKDIR /app

EXPOSE 8090

CMD ["./new_go"]
