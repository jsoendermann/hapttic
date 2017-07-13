FROM golang:alpine as build-env

MAINTAINER Jan Soendermann <jan.soendermann+git@gmail.com>

WORKDIR /hapttic-src
COPY . .

RUN go build -o hapttic .


FROM alpine

WORKDIR /

COPY --from=build-env /hapttic-src/hapttic .

RUN apk add --no-cache \
  bash \
  jq \
  curl \
  && rm -rf /var/cache/apk/*

EXPOSE 8080

ENTRYPOINT ["/hapttic"]
CMD ["-help"]
