# nsscache-go

[![GoDoc](https://godoc.org/github/MiLk/nsscache-go?status.png)](https://godoc.org/github/MiLk/nsscache-go)
[![Build Status](https://travis-ci.org/MiLk/nsscache-go.svg?branch=master)](https://travis-ci.org/MiLk/nsscache-go)
[![Coverage Status](https://coveralls.io/repos/github/MiLk/nsscache-go/badge.svg?branch=master)](https://coveralls.io/github/MiLk/nsscache-go?branch=master)

Implementation of [nsscache](https://github.com/google/nsscache) in Go.
The main goal of this library is too allow to write easily new program which can populate the nsscache files
from not yet supported sources or to use your custom logic to generate those cache files.

## Running the test

To run the test against [libnss-cache](https://github.com/google/libnss-cache),
you need to have docker installer, and build the test image.

```bash
cd docker
docker build -t nsscache-go:latest .
```

## See

* [nsscache](https://github.com/google/nsscache)
* [libnss-cache](https://github.com/google/libnss-cache)
