FROM debian:stable-slim

# COPY source destination
COPY goserver /bin/goserver

COPY .env /.env

ENV PORT=8080

CMD ["/bin/goserver"]
