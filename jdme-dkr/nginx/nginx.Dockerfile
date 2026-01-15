FROM alpine:3.19 AS imgopt
RUN apk add --no-cache oxipng jpegoptim
WORKDIR /img
COPY www/img/ ./
RUN oxipng -o2 --strip safe **/*.png \
 && jpegoptim --strip-all **/*.jpg **/*.jpeg

FROM alpine:latest AS fetch

RUN apk add --no-cache git 

ARG REPO_URL=https://github.com/jdetok/resume.git
ARG REPO_REF=main

RUN git clone --depth 1 --branch ${REPO_REF} ${REPO_URL} /site

FROM node:20-alpine AS tsbuild

RUN mkdir -p www/js/bball

WORKDIR /src

COPY package.json package-lock.json tsconfig.json ./
RUN npm ci

COPY ts ./ts

RUN npm run build

FROM nginx:stable-alpine

# # auth for private pages
COPY jdme-dkr/nginx/.htpasswd /etc/nginx/.htpasswd
RUN chmod 644 /etc/nginx/.htpasswd

RUN apk add --no-cache curl \
 && mkdir -p /etc/nginx /var/www /var/resume \
 && touch /var/log/nginx/nginx.log /var/log/nginx/err.log \
 && chmod 644 /etc/nginx/.htpasswd

# COPY www /var/www
COPY --from=fetch /site/public /var/resume
COPY --from=tsbuild /src/www/js /var/www/js
COPY --from=imgopt /img /var/www/img

RUN touch /var/log/nginx/nginx.log /var/log/nginx/err.log \
 && chown -R nginx:nginx /var/log/nginx /var/www /var/resume

HEALTHCHECK --interval=15s --timeout=5s --start-period=2s --retries=10 \
  CMD ["curl", "-f", "http://localhost/proxy-health"]
