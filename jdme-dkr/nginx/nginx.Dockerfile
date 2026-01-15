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
COPY www/js ./www/js

RUN npm run build

FROM alpine:3.19 AS imgopt

RUN apk add --no-cache oxipng jpegoptim

WORKDIR /img 

# copy images from project root
COPY www/img/pets ./pets 
COPY www/img/abt ./abt 
COPY www/img/slu_leds ./slu_leds 
COPY www/img/bronto ./bronto 
COPY www/img/*.* ./ 

# enable cache mount in build command:
# docker build --progress=plain --build-arg BUILDKIT_INLINE_CACHE=1 ...
RUN --mount=type=cache,target=/img_cache \
    find . -type f -name '*.png' -exec oxipng -o4 --strip safe {} + \
    && find . -type f \( -name '*.jpg' -o -name '*.jpeg' \) -exec jpegoptim --strip-all {} +


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
COPY --from=fetch /site/public /var/resume
COPY --from=tsbuild /src/www/js /var/www/js
COPY --from=imgopt /img /var/www/img


# create empty log files
RUN touch /var/log/nginx/nginx.log
RUN touch /var/log/nginx/err.log

HEALTHCHECK --interval=15s --timeout=5s --start-period=2s --retries=10 \
  CMD ["curl", "-f", "http://localhost/proxy-health"]
