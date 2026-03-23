FROM alpine:3.19 AS imgopt
RUN apk add --no-cache oxipng jpegoptim
WORKDIR /img
COPY www/img/ ./
RUN oxipng -o2 --strip safe **/*.png \
 && jpegoptim --strip-all **/*.jpg **/*.jpeg

FROM node:20-alpine AS tsbuild

RUN mkdir -p www/js/bball

WORKDIR /src

COPY package.json package-lock.json tsconfig.json ./
RUN npm ci

COPY ts ./ts

RUN npm run build

FROM node:current-alpine AS mrpbuild

WORKDIR /app

COPY ./stl-transit/package.json ./stl-transit/tsconfig.json ./stl-transit/vite.config.ts ./

RUN npm i

ENV NODE_OPTIONS="--max-old-space-size=4096"

COPY ./stl-transit/www/*.html www/
COPY ./stl-transit/www/css/*.css www/css/
COPY ./stl-transit/www/src/ www/src/

RUN npm run build

FROM nginx:stable-alpine

# auth for private pages
COPY dkr/nginx/.htpasswd /etc/nginx/.htpasswd
RUN chmod 644 /etc/nginx/.htpasswd

RUN apk add --no-cache curl \
 && mkdir -p /etc/nginx /var/www /var/resume /var/mrp \
 && touch /var/log/nginx/nginx.log /var/log/nginx/err.log \
 && chmod 644 /etc/nginx/.htpasswd

COPY --from=tsbuild /src/www/js /var/www/js
COPY --from=mrpbuild /app/www/js /var/mrp/js
COPY --from=imgopt /img /var/www/img

RUN touch /var/log/nginx/nginx.log /var/log/nginx/err.log \
 && chown -R nginx:nginx /var/log/nginx /var/www /var/resume

HEALTHCHECK --interval=15s --timeout=5s --start-period=2s --retries=10 \
  CMD ["curl", "-f", "http://localhost/proxy-health"]
