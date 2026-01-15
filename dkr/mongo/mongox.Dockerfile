FROM mongo-express:latest
USER root

RUN apk add --no-cache curl

EXPOSE 8081

HEALTHCHECK --interval=15s --timeout=5s --start-period=5s --retries=5 \
  CMD ["curl", "-f", "http://localhost:8081/status"]