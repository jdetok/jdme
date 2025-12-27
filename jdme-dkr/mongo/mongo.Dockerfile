FROM mongo:latest

COPY init/ /docker-entrypoint-initdb.d/

EXPOSE 27017

HEALTHCHECK --interval=10s --timeout=5s --start-period=20s --retries=5 \
  CMD ["mongosh", "--quiet", "--eval", "quit(db.adminCommand('ping').ok ? 0 : 1)"]
