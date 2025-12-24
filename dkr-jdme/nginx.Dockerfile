FROM nginx:alpine

# Copy config and htpasswd
COPY proxy/nginx.conf /etc/nginx/nginx.conf
COPY proxy/.htpasswd /etc/nginx/.htpasswd
# COPY proxy/ssl/cloudflare-origin.pem /etc/nginx/ssl/cloudflare-origin.pem
# COPY proxy/ssl/cloudflare-origin.key /etc/nginx/ssl/cloudflare-origin.key