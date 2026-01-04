FROM alpine:latest AS fetch

RUN apk add --no-cache git 

ARG REPO_URL=https://github.com/jdetok/resume.git
ARG REPO_REF=main

RUN git clone --depth 1 --branch ${REPO_REF} ${REPO_URL} /site

FROM nginx:stable

RUN touch /var/log/nginx/nginx.log
RUN touch /var/log/nginx/err.log

RUN mkdir -p /etc/nginx /var/www

# main nginx config
COPY nginx.conf /etc/nginx/nginx.conf

# auth for private pages
COPY .htpasswd /etc/nginx/.htpasswd
RUN chmod 644 /etc/nginx/.htpasswd

COPY --from=fetch /site/public /var/resume

# create empty log files
RUN touch /var/log/nginx/nginx.log
RUN touch /var/log/nginx/err.log

HEALTHCHECK --interval=15s --timeout=5s --start-period=5s --retries=10 \
  CMD ["curl", "-f", "http://localhost/proxy-health"]
