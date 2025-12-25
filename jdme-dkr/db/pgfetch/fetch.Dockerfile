FROM alpine:latest
# RUN apk add --no-cache postgresql-client bash
RUN apk add --no-cache bash tzdata
WORKDIR /fetch

# COPY pgbkp/pg_dump.sh /usr/local/bin/pg_dump.sh
COPY pgfetch/fetch.sh ./fetch.sh
RUN chmod +x ./fetch.sh

RUN echo "27 10 * * * /fetch/fetch.sh" > /etc/crontabs/root

# run cron in the foreground
CMD ["crond", "-f", "-L", "/var/log/cron.log"]