VERSION ?= v$(shell cat ./cmd/noop-server/version.go | grep Version | awk -F'"' '{print $$2}')
SHA ?= $(shell git reflog | head -n1 | cut -d' ' -f1)

default: test build

.PHONY: test
test:
	go clean -testcache
	go test -race ./...

.PHONY: benchmark bench
benchmark:
	# This target requires that 'tee' is installed
	go test -count=1 -benchmem -run='^$$' -count=1 -bench '^Benchmark.*$$' github.com/jmervine/noop-server/... | tee BENCHMARK.txt

bench: benchmark

.PHONY: docker
docker:
	docker build -t jmervine/noop-server:latest .
	docker tag jmervine/noop-server:latest jmervine/noop-server:$(HEX)
	docker tag jmervine/noop-server:latest jmervine/noop-server:$(VERSION)

.PHONY: release
release: clean test git/tag docker docker/push

.PHONY: push
docker/push: require_owner
	docker push jmervine/noop-server:latest
	docker push jmervine/noop-server:$(HEX)
	docker push jmervine/noop-server:$(VERSION)

.PHONY: push
run:
	go run ./cmd/noop-server/ --record

bin/noop-server:
	go build -o bin/noop-server ./cmd/noop-server/...

# Should run 'make build push' first.
.PHONY: tag
git/tag: require_owner
	git pull --tags
	git tag -f $(VERSION)
	git push --tags

# This is ugly and fragile, but meh, it'll work for now.
.PHONY: require_owner
require_owner:
	@# This is ugly and fragile, but meh, it'll work for now.
	@if ! [ "$(USER)" = "jmervine" ]; then echo "You are not the owner and cannot perform this task."; exit 1; fi

.PHONY: clean
clean:
	rm -rf bin

.PHONY: todos
todos:
	@git grep -n TODO | grep -v Makefile | awk -F':' '{ print " - TODO["$$1":"$$2"]:"$$NF }'

.PHONY: version
version:
	echo $(VERSION)