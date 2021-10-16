
FROM fedora:latest
RUN mkdir /srv/e212 && yum -y install sqlite
WORKDIR /srv/e212
COPY e212 .
COPY e212.sql .
COPY e212_cmd .
COPY templates/ templates/
COPY public/ public/
RUN ./e212_cmd newuser admin admin@example.com admin 
#sample date contains duplicates, makes sqlite3 fail
RUN sqlite3  mccmnc.db < e212.sql || true
EXPOSE 8080
ENTRYPOINT /srv/e212/e212 -port 8080
