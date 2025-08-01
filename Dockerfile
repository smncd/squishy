FROM golang:bullseye AS builder

COPY . /code

WORKDIR /code

RUN make build-linux-amd64

FROM alpine:3.22

WORKDIR /

COPY --from=builder /code/bin/* /squishy

RUN chmod +x /squishy

EXPOSE 1394

CMD [ "/squishy" ]
