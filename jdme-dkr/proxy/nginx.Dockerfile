FROM nginx:stable

RUN mkdir -p /etc/nginx/ssl 

RUN touch /var/log/nginx/nginx.log
RUN touch /var/log/nginx/err.log

COPY nginx.conf /etc/nginx/nginx.conf

COPY .htpasswd /etc/nginx/.htpasswd
RUN chmod 644 /etc/nginx/.htpasswd

# COPY ssl/cloudflare-origin.pem /etc/nginx/ssl/cloudflare-origin.pem
# COPY ssl/cloudflare-origin.key /etc/nginx/ssl/cloudflare-origin.key

HEALTHCHECK --interval=15s --timeout=5s --start-period=15s --retries=10 \
  CMD ["curl", "-f", "https://jdeko.me/proxy-health"]
