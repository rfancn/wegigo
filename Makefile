.PHONY: shared
all: shared

BINARY=wegigo

dev:
	go generate
	CGO_ENABLED=1 go build -ldflags -s -o $(BINARY)
static: go_dep
	CGO_ENABLED=0 go build -ldflags -s -o $(BINARY)
shared: go_dep
	CGO_ENABLED=1  go build -ldflags -s -o $(BINARY)
test: go_dep
	CGO_ENABLED=0 go build -ldflags="-X deploy.MODE dev" -ldflags -s -o wegigo

# ensure go library dependencies
go_dep:
	go get -u -v github.com/golang/dep/cmd/dep
	go get -u github.com/jteeuwen/go-bindata/...
	dep ensure
	go generate

ansible_dep:
	git

# build plugin
APP_DIR ?= apps/$(APP_NAME)
plugin:
	# generate app plugin asset data before go build
	go-bindata -o $(APP_DIR)/bindata.go $(APP_DIR)/asset/...
	CGO_ENABLED=1 go build -buildmode=plugin -o $(APP_DIR)/$(APP_NAME).so $(APP_DIR)/*.go




