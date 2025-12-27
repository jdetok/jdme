FROM golang:1.25

WORKDIR /app

# RUN go install github.com/air-verse/air@latest
COPY go.mod go.sum ./

RUN go mod download

COPY . .

# arm64 for prod (pi) amd64 for mac 
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/api ./main

ENTRYPOINT [ "/app/bin/api" ]

HEALTHCHECK --interval=15s --timeout=5s --start-period=15s --retries=10 \
  CMD ["curl", "-f", "http://localhost:8080/health"]
  # CMD ["curl", "-f", "http://localhost:8080/health"]