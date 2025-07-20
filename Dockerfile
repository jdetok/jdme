FROM golang:1.24

# testing 07/20/2025 with refactored frontend in separate dir
# RUN mkdir -p /static
# COPY /home/jdeto/frontend_jdme/. /static/.

WORKDIR /app
RUN go install github.com/air-verse/air@latest
COPY . .
RUN go mod download

ENTRYPOINT [ "air" ]
