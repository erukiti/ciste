build:
	go fmt *.go
	go build

dist-update:
	git submodule foreach 'git checkout master; git pull'
	go-bindata -pkg=main -o=dist.go ./ciste-web-content/dist/*
	go fmt *.go

init:
	git submodule update --init --recursive

get:
	go get
	go get -u github.com/jteeuwen/go-bindata/...
