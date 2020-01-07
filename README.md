# go-maxconn
go-maxconn is a small tool that opens many TCP connections to a server for testing

## How to build

Easy:
```
go build
```

## How to use
```
Usage of ./go-maxconn:
  -address string
    	the address to connect to e.g. www.myhost.com:443
  -connections int
    	maximum number of concurrent connections to open (default 100)
  -tls
    	whether to use TLS to connect or not (default true)
  -wait duration
    	time to wait before tearing down connections again (default 5m0s)
```

## Examples
```
./go-maxconn -address myhost.com:443 -connections 100 -wait 1m
```

This will open 100 TLS connections to `myhost.com` and close them again after one minute.

```
./go-maxconn -address localhost:8080 -tls=false -connections 10 -wait 1m
```

This will open 10 TCP connections to `localhost` and close them after one minute.
(notice that you have to pass `-tls=false` to make it work. `-tls false` wont' work)
