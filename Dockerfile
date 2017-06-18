FROM golang:latest

MAINTAINER Jan Soendermann <jan.soendermann+git@gmail.com>

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app
COPY . .

RUN go build -o hapttic .

EXPOSE 8080

CMD ["/usr/src/app/hapttic", "-f", "/hapttic_request_handler.sh"]
