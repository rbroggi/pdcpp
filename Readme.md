# Pdc price provider module

This module acts like an in-memory cache of prices.  
The output tool is a command-line tool designed using cobra that enables users to serve grpc-based API's and cosume them with built-in clients.

## Getting Started

pdcpp [subcomand] [flags] 

```
$ pdcpp serve --port 8080 &
$ pdcpp ping localhost:8080
```

For TLS authentication refer to certs/Makefile and the following blog post: [bbengfort](https://bbengfort.github.io/programmer/2017/03/03/secure-grpc.html)

### Prerequisites

* [Golang](https://golang.org/) - The language 
* [GRPC](https://grpc.io/) - Communication protocol 
* [Cobra](https://github.com/spf13/cobra) - Command line tool framework

### Installing

```
$ make goget  
$ make  
```

## Running the tests

```
$ make test
```

## Deployment

The make commad will output a single simple binary into the *target* directory - you can then copy it to the target run environment provided that it is architecturally compatible with the building environment. (Please check golang manual for [cross-compilation](https://golang.org/doc/install/source#environment) necessity)

## Built With

* [Golang](https://golang.org/) - The language 
* [GRPC](https://grpc.io/) - Communication protocol 
* [Cobra](https://github.com/spf13/cobra) - Command line tool framework

## Authors

* **Rodrigo Broggi** - *Initial work* - [rbroggi](https://github.com/rbroggi)

