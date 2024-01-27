FROM golang:1.21-alpine3.19 AS builder

COPY . /GRPC_Family/

WORKDIR /GRPC_Family/

RUN go mod download

RUN go build -o ./bin/app cmd/family/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 GRPC_Family/bin/app .
COPY --from=0 GRPC_Family/config config/

EXPOSE 80

CMD ["./app"]