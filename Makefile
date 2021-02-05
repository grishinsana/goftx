test:
	go test -v ../.

lint:
	CGO_ENABLED=1 golangci-lint run --config=./.golangci.yml
