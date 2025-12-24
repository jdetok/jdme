FROM postgres:17
USER postgres

WORKDIR /docker-entrypoint-initdb.d
COPY ./db/dump/. .

WORKDIR /var/lib/postgresql
RUN mkdir dump

COPY ./db/dump/. dump/.