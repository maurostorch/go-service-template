FROM golang:1.18 as builder

WORKDIR /opt/build
COPY Makefile .
COPY go.mod .
COPY go.sum .
RUN make dep

COPY main.go .
COPY handler ./handler

RUN make test build

FROM alpine
RUN apk add ca-certificates
COPY --from=builder /opt/build/app /bin/app

CMD [ "/bin/app" ]
