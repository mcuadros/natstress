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
  --shutdown    Wait time for received all the messages sent. (Default: 5s)
```

Output
------

`./natstress -shutdown=10s -m=10000 -r=10000 -c=500`

```
Creating 500 client(s) ... OK
Subscribing clients to 5 subject(s) ... OK
Publishing and receiving messages...
Received 5000000 / 5000000 [================================] 100.00 % 611714/s

Punlishing summary:
  Count:    10000 messages.
  Total:    1.3490 secs.
  Slowest:  3245 µs.
  Fastest:  0002 µs.
  Average:  0007 µs.
  Messages/sec: 7412.9942
``

License
-------

MIT, see [LICENSE](LICENSE)
