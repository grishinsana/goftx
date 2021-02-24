checkmaster:
	@if [ -z ${BRANCH} ] || [ ${BRANCH} != "master" ]; then \
		git checkout master && git pull; \
    fi

tagger: checkmaster tag-m

tag-m:
	./scripts/tag.sh --minor

tag-p:
	./scripts/tag.sh --patch

test:
	go test -v ../.

lint:
	CGO_ENABLED=1 golangci-lint run --config=./.golangci.yml
