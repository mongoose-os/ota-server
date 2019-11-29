# Mongoose OS OTA firmware server

This server serves firmware for Mongoose OS devices.

While Mongoose OS firmware can be served by any HTTP server, this implementation provides two distinct features:

 * Optimizes behavior when version does not change. If the version sent by the device matches the version in the ZIP file, "304 Not  modified is returned" and no update is performed.
 * Memory usage in TLS mode is optimized by forcing smaller TLS messages (1K).

## Building

```
$ go build -v
```

## Usage

### HTTP only
```
$ ./ota-server --root=/var/www/html --listen-addr=:8910 --alsologtostderr
```

### HTTP and HTTPS
```
$ ./ota-server --root=/var/www/html --listen-addr=:8910 --tls-listen-addr=:8443 --tls-cert=cert.pem --tls-key=key.pem --alsologtostderr

```
