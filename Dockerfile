FROM debian:stable-slim

# COPY source destination
COPY goserver /bin/goserver

COPY .env /.env


CMD ["/bin/goserver"]
