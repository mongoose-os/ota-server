# Mongoose OS OTA firmware server

This server serves firmware for Mongoose OS devices.

While Mongoose OS firmware can be served by any HTTP server, this implementation provides two distinct features:

 * Optimizes behavior when version does not change. If the version sent by the device matches the version in the ZIP file, "304 Not  modified is returned" and no update is performed.
 * Memory usage in TLS mode is optimized by forcing smaller TLS messages (1K).

## Building

```
$ go build -v
```

## Docker image

There is a public Docker image hosted on Docker Hub, [docker.io/mgos/ota-server](https://hub.docker.com/r/mgos/ota-server):

```
$ docker pull docker.io/mgos/ota-server:latest
```

## Usage

### HTTP only
```
$ ./ota-server --root=/var/www/html --listen-addr=:8910 --alsologtostderr
```

### HTTP and HTTPS
```
$ ./ota-server --root=/var/www/html --listen-addr=:8910 --tls-listen-addr=:8443 \
               --tls-cert=cert.pem --tls-key=key.pem --alsologtostderr

```

### Using Docker image

Persistent container with [LetsEncrypt](https://letsencrypt.org/) certificates.

```
$ docker run --name ota-server --detach --restart always \
             -v /etc/letsencrypt:/etc/letsencrypt -v /home/ubuntu/www:/www \
             -p 80:8080 -p 443:8443 docker.io/mgos/ota-server:latest /ota-server \
                --root=/www --listen-addr=:8080 --tls-listen-addr=:8443 \
                --tls-cert=/etc/letsencrypt/live/MY-DOMAIN.com/fullchain.pem \
                --tls-key=/etc/letsencrypt/live/MY-DOMAIN.com/privkey.pem
```
