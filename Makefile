.PHONY: static
all: static

BINARY=wegigo

dev: go_dep
	CGO_ENABLED=0 go build -ldflags="-X deploy.MODE dev" -ldflags -s -o wegigo
static: go_dep
	CGO_ENABLED=0 go build -ldflags -s -o $(BINARY)
shared: go_dep
	CGO_ENABLED=1  go build -ldflags -s -o wegigo

# ensure go library dependencies
go_dep:
	go get -u -v github.com/golang/dep/cmd/dep
	go get -u github.com/jteeuwen/go-bindata/...
	dep ensure
	go generate

ansible_dep:
	git

# build plugin
plugin:
	CGO_ENABLED=0 go build -buildmode=plugin  -o plugins/mod_deploy.so pkg/deploy/module.go
