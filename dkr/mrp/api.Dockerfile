FROM golang:1.26-alpine

RUN apk add --no-cache git curl

ARG REPO_URL=https://github.com/jdetok/stl-transit.git
ARG REPO_REF=main

RUN git clone --depth 1 --branch ${REPO_REF} ${REPO_URL} /app

WORKDIR /app

COPY .env ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bin/app ./src

ENTRYPOINT [ "/app/bin/app" ]

HEALTHCHECK --interval=5s --timeout=3s --start-period=10s --retries=10 \
    CMD [ "curl", "-f", "http://localhost:9999/health" ]