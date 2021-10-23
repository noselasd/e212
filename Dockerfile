FROM golang:latest AS builder
WORKDIR /go/src/app
COPY . .
RUN make

FROM fedora:latest
RUN mkdir /srv/e212 && yum -y install sqlite
WORKDIR /srv/e212
COPY --from=builder /go/src/app/e212 .
COPY --from=builder /go/src/app/e212_init_data.sh .
COPY --from=builder /go/src/app/e212_cmd .
COPY --from=builder /go/src/app/templates/ templates/
COPY --from=builder /go/src/app/public/ public/
RUN ./e212_cmd newuser admin admin@example.com admin 
#sample date contains duplicates, makes sqlite3 fail
RUN ./e212_init_data.sh
EXPOSE 8080
ENTRYPOINT /srv/e212/e212 -port 8080
