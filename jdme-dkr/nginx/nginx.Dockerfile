FROM alpine:latest AS fetch

RUN apk add --no-cache git 

ARG REPO_URL=https://github.com/jdetok/resume.git
ARG REPO_REF=main

RUN git clone --depth 1 --branch ${REPO_REF} ${REPO_URL} /site

FROM node:20-alpine AS tsbuild

WORKDIR /src

COPY package.json package-lock.json tsconfig.json ./
RUN npm ci

COPY ts ./ts
COPY www ./www

RUN npm run build

FROM nginx:stable

RUN touch /var/log/nginx/nginx.log
RUN touch /var/log/nginx/err.log

RUN mkdir -p /etc/nginx /var/www

# main nginx config
# COPY jdme-dkr/nginx/nginx.conf /etc/nginx/nginx.conf

# auth for private pages
COPY jdme-dkr/nginx/.htpasswd /etc/nginx/.htpasswd
RUN chmod 644 /etc/nginx/.htpasswd

# COPY www /var/www

COPY --from=tsbuild /src/www /var/www
COPY --from=fetch /site/public /var/resume

# create empty log files
RUN touch /var/log/nginx/nginx.log
RUN touch /var/log/nginx/err.log

HEALTHCHECK --interval=15s --timeout=5s --start-period=2s --retries=10 \
  CMD ["curl", "-f", "http://localhost/proxy-health"]
