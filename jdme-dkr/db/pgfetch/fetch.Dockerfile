#BUILD STAGE
FROM golang:1.25 AS builder

RUN apt-get update && apt-get install -y git \
    && rm -rf /var/lib/apt/lists/*

ARG REPO_URL=https://github.com/jdetok/bball-etl-cli.git
ARG REPO_REF=main

RUN git clone --depth 1 --branch ${REPO_REF} ${REPO_URL} /app && ls -al /app

# RUN ls /app

WORKDIR /app

RUN ls -al
RUN ls -al main
RUN mkdir ./log
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/fetch ./main

# RUNTIME STAGE
FROM alpine:latest

RUN apk add --no-cache bash postgresql-client tzdata

RUN mkdir /dump

COPY --from=builder /app /app

COPY ./pgfetch/fetch.sh /app/fetch.sh
RUN chmod +x /app/fetch.sh

RUN echo "35 00 * * * /bin/sh -c 'cd /app && ./fetch.sh'" > /etc/crontabs/root

CMD ["crond", "-f", "-L", "/var/log/cron.log"]

WORKDIR /app

HEALTHCHECK --interval=30s --timeout=5s --retries=5 \
    CMD pgrep crond >/dev/null || exit 1