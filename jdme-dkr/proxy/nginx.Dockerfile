FROM nginx:stable

RUN mkdir -p /etc/nginx/ssl /var/www

# main nginx config
COPY nginx.conf /etc/nginx/nginx.conf

# auth for private pages
COPY .htpasswd /etc/nginx/.htpasswd
RUN chmod 644 /etc/nginx/.htpasswd

# ssl certs
# COPY ssl/cloudflare-origin.pem /etc/nginx/ssl/cloudflare-origin.pem
# COPY ssl/cloudflare-origin.key /etc/nginx/ssl/cloudflare-origin.key

# create empty log files
RUN touch /var/log/nginx/nginx.log
RUN touch /var/log/nginx/err.log

HEALTHCHECK --interval=15s --timeout=5s --start-period=15s --retries=10 \
  CMD ["curl", "-f", "http://localhost/proxy-health"]
