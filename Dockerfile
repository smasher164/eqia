FROM golang:1.14-alpine
RUN apk --no-cache add ca-certificates
WORKDIR /eqia
COPY . .
RUN GOOS=linux go build -tags 'osusergo netgo'
FROM scratch
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=0 /eqia/eqia /
CMD ["/eqia"]