FROM golang:alpine

MAINTAINER Jan Soendermann <jan.soendermann+git@gmail.com>

RUN apk add --no-cache \
  bash \
  jq \
  curl \
  && rm -rf /var/cache/apk/*

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app
COPY . .

RUN go build -o hapttic .

EXPOSE 8080

ENTRYPOINT ["/usr/src/app/hapttic"]
CMD ["-help"]
