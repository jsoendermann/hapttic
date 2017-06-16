# Why would I want this?

- You want to run some code in response to a webhook, for example a github push. For an example, see [pfeife](/jsoendermann/pfeife)
- You have some code on your Raspberry Pi that you want to run from work (great in combination with [ngrok](https://ngrok.com/))
- You want to convert a webhook into an API call.

# How does it work?

Hapttic dumps all the information in 

- more details
- json format
- example script
- example docker compose
- how to install (in a dockerfile)
- will you add feature X?

# FAQ

### Does hapttic come with support for SSL?

You can add SSL to hapttic by putting an nginx proxy in front of it like so:

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
      - /var/letsencrypt_certs:/etc/nginx/certs:ro
      - vhost:/etc/nginx/vhost.d
      - html:/usr/share/nginx/html
    labels:
      - "com.github.jrcs.letsencrypt_nginx_proxy_companion.nginx_proxy=true"

  letsencrypt-nginx-proxy-companion:
    restart: always
    image: jrcs/letsencrypt-nginx-proxy-companion
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
      - /var/letsencrypt_certs:/etc/nginx/certs:rw
      - vhost:/etc/nginx/vhost.d
      - html:/usr/share/nginx/html

  hapttic:
    restat: always
    image: TODO
    environment:
      - VIRTUAL_HOST=hapttic.example.com                                                    # Replace this
      - LETSENCRYPT_HOST=hapttic.example.com                                                # Replace this
      - LETSENCRYPT_EMAIL=your@email.address                                                # Replace this
    depends_on:
      - nginx-proxy
      - letsencrypt-nginx-proxy-companion
```

### Will you add feature X?

Probably not, the primary goal of hapttic is ease of use. Check out shell2http or adnanh/webhook for more feature rich alternatives.

### Will this work on Windows?

Probably not because the path to "/bin/bash" is hardcoded. Pull requests welcome!

# Extra disclaimer

There are obvious security implications with running a bash script in response to . Be extra careful when you use this and have a look at the source (~100 lines of Go). Hapttic comes without any warranty, as per [the license](https://github.com/jsoendermann/hapttic/blob/master/LICENSE).
