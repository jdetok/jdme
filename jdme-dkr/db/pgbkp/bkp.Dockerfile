FROM alpine:latest
# RUN apk add --no-cache postgresql-client bash
RUN apk add --no-cache bash postgresql-client tzdata
WORKDIR /wrk

# COPY pgbkp/pg_dump.sh /usr/local/bin/pg_dump.sh
COPY pgbkp/pg_dump.sh ./pg_dump.sh
RUN chmod +x ./pg_dump.sh

RUN echo "27 10 * * * /wrk/pg_dump.sh" > /etc/crontabs/root

# run cron in the foreground
CMD ["crond", "-f", "-L", "/var/log/cron.log"]