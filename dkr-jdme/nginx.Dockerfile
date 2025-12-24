FROM nginx:stable

# Copy config and htpasswd

RUN mkdir -p /etc/nginx/ssl

COPY proxy/nginx.conf /etc/nginx/nginx.conf
COPY proxy/.htpasswd /etc/nginx/.htpasswd

RUN chmod 644 /etc/nginx/.htpasswd
# COPY proxy/ssl/cloudflare-origin.pem /etc/nginx/ssl/cloudflare-origin.pem
# COPY proxy/ssl/cloudflare-origin.key /etc/nginx/ssl/cloudflare-origin.key