FROM ubuntu:latest

WORKDIR /home

COPY ./dist/server ./server
COPY .env .env

RUN apt-get update && \
    apt-get install -y ca-certificates && \
    update-ca-certificates

RUN chmod +x /home/server

CMD ["/home/server"]
