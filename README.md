natstress [![Build Status](https://travis-ci.org/mcuadros/natstress.png?branch=master)](https://travis-ci.org/mcuadros/natstress) [![GoDoc](http://godoc.org/github.com/mcuadros/natstress?status.png)](http://godoc.org/github.com/mcuadros/natstress)
==============================

A stress tool for NATS Servers


Installation
------------

The recommended way to install natstress

```
go get github.com/mcuadros/natstress
```

A binnary called `natstress` should be at `$GOPATH/bin/`

Usage
-----

```go
Usage: natstress [options...]

Options:
  -h            NATS server url. (Default: nats://localhost:4222)
  -s            Number of subjects.
  -m            Number of message to send in each subject.
  -c            Number of clients to run concurrently.
  --warmup      Time to wait before start to deliver messages after connect
                to the server. (Default: 50ms)
  --shutdown    Wait time for received all the messages sent. (Default: 5s)
```

License
-------

MIT, see [LICENSE](LICENSE)
