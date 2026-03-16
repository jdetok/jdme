FROM postgres:18

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        postgis postgresql-18-postgis-3 osmium-tool osm2pgsql ca-certificates curl unzip gdal-bin \
    && rm -rf /var/lib/apt/lists/*

HEALTHCHECK --interval=5s --timeout=5s --start-period=5s --retries=10 \
  CMD ["pg_isready", "-h", "localhost", "-U", "gis"]