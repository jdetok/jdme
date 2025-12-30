FROM postgres:17

WORKDIR /var/lib/postgresql

# Create dump directory and copy files
RUN mkdir -p /var/lib/postgresql/dump
COPY dump/*.sql.gz /var/lib/postgresql/dump/

# sql scripts to create schemas, users, etc
COPY sql/a_init/. /docker-entrypoint-initdb.d/

# init.sh file for building 
COPY init.sh /docker-entrypoint-initdb.d/d_init.sh
RUN chmod +x /docker-entrypoint-initdb.d/d_init.sh

HEALTHCHECK --interval=10s --timeout=5s --start-period=20s --retries=5 \
  CMD ["pg_isready", "-h", "localhost", "-U", "postgres"]