# hapttic

## What is this good for?

- You want to run some code in response to a webhook, like a github push.
- You have some code on your Raspberry Pi that you want to run from work (great in combination with [ngrok](https://ngrok.com/)).
- That's pretty much it.

## How does it work?

Hapttic listens for incoming HTTP connections. When it receives a connection, it dumps all relevant data (headers, path, parameters, the body and other stuff) into a JSON object and calls a bash script with this object as its parameters.

## Isn't this just a reinvention of CGI?

The basic idea is pretty similar. The main advantage of hapttic is ease of setup with a simple Docker image that lets you quickly connect a shell script to a http endpoint.

## Show me an example

First, create a request handler at `~/hapttic_request_handler.sh`:

```bash
echo $1
```

Then run the following command to spin up the docker container that runs hapttic:

`docker run --rm -p 8080:8080 -v ~/hapttic_request_handler.sh:/hapttic_request_handler.sh --name hapttic jsoendermann/hapttic -file "/hapttic_request_handler.sh"`

Finally, run `open http://localhost:8080` to see the output of your script.

## Show me a more realistic example

```bash
REQUEST=$1
SECRET_TOKEN=$(jq -r '.Header."X-My-Secret"[0]' <(echo $REQUEST))

if [[ "$SECRET_TOKEN" != "SECRET" ]]; then
  echo "Incorrect secret token"
  exit -1
fi

curl https://www.example.com/api/call/in/response/to/webhook
```

This request handling script can be run with `curl -H "X-My-Secret: SECRET" http://localhost:8080`

The [`jsoendermann/hapttic`](https://hub.docker.com/r/jsoendermann/hapttic/) Dockerfile includes `jq` and `curl`, if you need any other command in your request handling script, you should create your own image.

## The Request JSON object

The JSON object your request handling script gets called with is a subset of Go's `http.Request`. It's defined in [hapttic.go](https://github.com/jsoendermann/hapttic/blob/master/hapttic.go) as `marshallableRequest`. For documentation on http.Request, see [the official net/http page](https://golang.org/pkg/net/http/#Request).

## SSL Support

You can add encryption by putting an nginx proxy in front of it with a docker-compose file like so:

```yaml
version: '3'

volumes:
  vhost:
  html:

services:
  nginx-proxy:
    restart: always
    image: jwilder/nginx-proxy
    ports:
      - 80:80
      - 443:443
    volumes:
      - /var/run/docker.sock:/tmp/docker.sock:ro
      - /var/certs:/etc/nginx/certs:ro
      - vhost:/etc/nginx/vhost.d
      - html:/usr/share/nginx/html
    labels:
      - "com.github.jrcs.letsencrypt_nginx_proxy_companion.nginx_proxy=true"

  letsencrypt-nginx-proxy-companion:
    restart: always
    image: jrcs/letsencrypt-nginx-proxy-companion
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
      - /var/certs:/etc/nginx/certs:rw
      - vhost:/etc/nginx/vhost.d
      - html:/usr/share/nginx/html

  hapttic:
    restart: always
    image: jsoendermann/hapttic
    environment:
      - VIRTUAL_HOST=hapttic.your.domain.com                                # Replace this
      - LETSENCRYPT_HOST=hapttic.your.domain.com                            # Replace this
      - LETSENCRYPT_EMAIL=your@email.address                                # Replace this
    volumes:
      - /my-request-handler.sh:/hapttic_request_handler.sh                  # Replace this
    command: ["-file", "/hapttic_request_handler.sh"]
    depends_on:
      - nginx-proxy
      - letsencrypt-nginx-proxy-companion
```
