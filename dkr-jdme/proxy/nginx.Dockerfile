FROM nginx:stable

RUN mkdir -p /etc/nginx/ssl

COPY nginx.conf /etc/nginx/nginx.conf
COPY .htpasswd /etc/nginx/.htpasswd

RUN chmod 644 /etc/nginx/.htpasswd
# COPY proxy/ssl/cloudflare-origin.pem /etc/nginx/ssl/cloudflare-origin.pem
# COPY proxy/ssl/cloudflare-origin.key /etc/nginx/ssl/cloudflare-origin.key

HEALTHCHECK --interval=15s --timeout=5s --start-period=15s --retries=10 \
  CMD ["curl", "-f", "http://localhost:8080/proxy-health"]