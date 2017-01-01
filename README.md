nildns
======

A minimum DNS for nginx + Docker.

# What's this?

This is a minimum DNS server which is designed for Nginx reverse proxy liked to Docker containers.
`nildns` uses < 2kB memory and 5 times less memory than `go-dnsmasq-min`.

#### Run

```
go run nildns.go -address=127.0.0.1:1153 -conf=./resolv.conf
```

#### Options

```
Usage of nildns
  -address string
    Listen address (default "127.0.0.1:53")
  -conf string
    Path to resolv.conf (default "/etc/resolv.conf")
  -tcp
    Enable TCP
  -ttl int
    Default TTL (default 10)
```
