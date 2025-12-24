FROM mongo:latest

COPY mongo/init/ /docker-entrypoint-initdb.d/

EXPOSE 27017