.PHONY: vet linux osx build test release

vet:
	go tool vet *.go

linux:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/terraform-provider-scaleft-linux .

osx:
	GOOS=darwin GOARCH=386 go build -o bin/terraform-provider-scaleft-osx .

build: vet osx linux
	go install .

test: build
	go test -v ./...

release:
	./bin/release.sh
