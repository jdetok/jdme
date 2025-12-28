FROM mongo-express:latest
USER root

RUN apk add --no-cache curl

EXPOSE 8081

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=5 \
  CMD ["curl", "-f", "http://localhost/status"]