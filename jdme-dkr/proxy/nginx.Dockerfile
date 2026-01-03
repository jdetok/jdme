FROM nginx:stable

# RUN mkdir -p /etc/nginx/ssl 

RUN touch /var/log/nginx/nginx.log
RUN touch /var/log/nginx/err.log

COPY nginx.conf /etc/nginx/nginx.conf

COPY .htpasswd /etc/nginx/.htpasswd
RUN chmod 644 /etc/nginx/.htpasswd

HEALTHCHECK --interval=15s --timeout=5s --start-period=15s --retries=10 \
  CMD ["curl", "-f", "http://localhost/proxy-health"]