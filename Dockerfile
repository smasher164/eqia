FROM golang:1.14
COPY . /eqia/
WORKDIR /eqia
RUN GOOS=linux go build -tags 'osusergo netgo'
FROM scratch
COPY --from=0 /eqia/eqia /
CMD ["/eqia"]