
FROM fedora:latest
RUN mkdir /srv/e212
WORKDIR /srv/e212
COPY e212 .
COPY e212_cmd .
COPY templates/ .
RUN ./e212_cmd newuser admin admin@example.com admin
EXPOSE 8080
ENTRYPOINT /srv/e212/e212 -port 8080
