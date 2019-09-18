FROM golang:1.12 AS builder
COPY . /build/
WORKDIR /build/
RUN go build -mod=vendor .

FROM alpine:3.10
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
WORKDIR /pta/
COPY --from=builder /build/app /pta/
COPY --from=builder /build/aux-addon.lua /pta/
EXPOSE 8111
ENTRYPOINT sh -c "sleep 10 && ./app"
