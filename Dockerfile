FROM --platform=$BUILDPLATFORM golang:bullseye AS builder

ARG TARGETOS

ARG TARGETARCH

WORKDIR /code

COPY go.mod go.sum /code/

RUN go mod download

COPY . .

RUN make build-${TARGETOS}-${TARGETARCH}


FROM alpine:3.22

WORKDIR /

COPY --from=builder /code/bin/* /squishy

RUN chmod +x /squishy

EXPOSE 1394

CMD [ "/squishy" ]
