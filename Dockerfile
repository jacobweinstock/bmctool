FROM golang:1.16 as builder

WORKDIR /code
COPY go.mod go.sum /code/
RUN go mod download

COPY . /code
RUN make build

FROM alpine:3.8
LABEL maintainers="https://github.com/jacobweinstock"

WORKDIR /tmp
RUN apk add ipmitool=1.8.18-r9

USER nobody
COPY --from=builder /code/bin/bmctool-linux-amd64 /bmctool-linux-amd64

ENTRYPOINT ["/bmctool-linux-amd64"]