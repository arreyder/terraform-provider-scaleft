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
	TF_ACC=yes MESOS_KAFKA_URL="http://dev.banno.com:7000" go test -v ./...

release:
	./bin/release.sh
