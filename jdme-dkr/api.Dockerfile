FROM golang:1.25

WORKDIR /app

# RUN go install github.com/air-verse/air@latest
COPY go.mod go.sum ./

RUN go mod download

COPY .env .env
COPY ./api ./api
COPY ./pkg ./pkg
COPY ./main ./main
COPY ./persist ./persist

# arm64 for prod (pi) amd64 for mac 
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bin/api ./main

ENTRYPOINT [ "/app/bin/api" ]

HEALTHCHECK --interval=15s --timeout=5s --start-period=5s --retries=10 \
  CMD ["curl", "-f", "http://localhost:8080/health"]