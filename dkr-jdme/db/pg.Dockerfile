FROM postgres:17
USER postgres

WORKDIR /docker-entrypoint-initdb.d
COPY dump/. .

WORKDIR /var/lib/postgresql
RUN mkdir dump

COPY dump/. dump/.

HEALTHCHECK --interval=10s --timeout=5s --start-period=20s --retries=5 \
  CMD ["pg_isready", "-h", "localhost", "-U", "postgres"]

