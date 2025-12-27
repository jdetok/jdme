FROM postgres:17

WORKDIR /var/lib/postgresql

COPY init.sh /usr/local/bin/init.sh
RUN chmod +x /usr/local/bin/init.sh

RUN mkdir dump

COPY dump/prod.sql dump/.

# ENTRYPOINT ["/usr/local/bin/init.sh"]
# CMD ["postgres"]

HEALTHCHECK --interval=10s --timeout=5s --start-period=20s --retries=5 \
  CMD ["pg_isready", "-h", "localhost", "-U", "postgres"]

