Docker registry pruner
======================

This is a tool for batch deleting manifests from a private Docker registry.

`docker-registry-pruner` [deletes through the registry HTTP API](https://docs.docker.com/registry/spec/api/#deleting-an-image)
so it does not actually remove images from the registry storage, is just removes the manifest and tags. If you then want
to free up some space, [run the registry garbage collector](https://docs.docker.com/registry/garbage-collection/#run-garbage-collection)
which will remove dangling layers from storage.

Configure and run it
--------------------

`docker-registry-pruner` accepts parameters as flags (eg `./docker-registry-pruner -minage 24h`) and as environment
variables with the prefix `DOCKER_REGISTRY_PRUNER_` (eg. `DOCKER_REGISTRY_PRUNER_MINAGE=24h ./docker-registry-pruner`).

To get a description of the available flags, run `./docker-registry-pruner -h`.

Building
--------

`docker-registry-pruner` is written in Go and can be build with the usual `go get . && go build .` from within your `$GOPATH`.
If you'd rather not set up a Go build environment on your machine, you can also built it using the golang docker image, like so:

    docker run --rm \
               -v /tmp/golang_docker-registry-pruner:/go \
               -v $(pwd):/go/src/github.com/cego/docker-registry-pruner \
               -w /go/src/github.com/cego/docker-registry-pruner \
               golang:1.8 go get .
    docker run --rm \
               -v /tmp/golang_docker-registry-pruner:/go \
               -v $(pwd):/go/src/github.com/cego/docker-registry-pruner \
               -w /go/src/github.com/cego/docker-registry-pruner \
               -e CGO_ENABLED=0 \
               golang:1.8 go build -a .

