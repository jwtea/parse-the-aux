FROM golang:1.12 AS builder

COPY . ./app
RUN cd ./app && pwd
RUN cd ./app && \
  go get . && \
  go install github.com/jwtea/parse-the-aux/app; 

FROM alpine:3.7  

COPY --from=builder /go/bin/app .

EXPOSE 8111

CMD ["sleep", "3000"]