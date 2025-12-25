FROM postgres:17
USER postgres

WORKDIR /docker-entrypoint-initdb.d
COPY ./dump/. .

WORKDIR /var/lib/postgresql
RUN mkdir ./dump

COPY ./dump/. ./dump/.

# RUN mkdir -p /var/lib/postgresql/dump
# RUN mkdir -p /wrk

# COPY init.sh /wrk/init.sh
# RUN chmod +x /wrk/init.sh

# COPY init.sh /usr/local/bin/init.sh
# RUN chmod +x /usr/local/bin/init.sh
# ENTRYPOINT ["/usr/local/bin/init.sh"]
# CMD ["postgres"]

HEALTHCHECK --interval=10s --timeout=5s --start-period=20s --retries=5 \
  CMD ["pg_isready", "-h", "localhost", "-U", "postgres"]

