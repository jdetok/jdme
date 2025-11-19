FROM httpd:2.4

WORKDIR /usr/local/apache2

# COPY SSL CERT AND KEY
COPY ssl/cloudflare-origin.pem ssl/cloudflare-origin.pem
COPY ssl/cloudflare-origin.key ssl/cloudflare-origin.key
COPY jdeko.conf conf/extra/jdeko.conf
COPY httpd.conf conf/httpd.conf

EXPOSE 80 443
