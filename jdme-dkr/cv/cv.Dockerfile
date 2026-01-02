FROM alpine:latest

RUN apk add --no-cache git 

ARG REPO_URL=https://github.com/jdetok/resume.git
ARG REPO_REF=main

RUN git clone --depth 1 --branch ${REPO_REF} ${REPO_URL} /site

WORKDIR /site