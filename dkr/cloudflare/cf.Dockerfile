FROM cloudflare/cloudflared:latest

WORKDIR /etc/cloudflared
COPY conf/ .