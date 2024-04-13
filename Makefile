default: test build

.PHONY: test
test:
	go test -v -race ./...

.PHONY: build
build:
	docker build -t jmervine/noop-server:latest .
	docker tag jmervine/noop-server:latest jmervine/noop-server:$(shell git reflog | head -n1 | cut -d' ' -f1)

.PHONY: push
push:
	docker push jmervine/noop-server:latest
	docker push jmervine/noop-server:$(shell git reflog | head -n1 | cut -d' ' -f1)

.PHONY: push
run:
	go run ./cmd/noop-server/

bin/noop-server:
	go build -o bin/noop-server ./cmd/noop-server/...

# Should run 'make build push' first.
.PHONY: tag
tag: require_tag
	# Tag git and docker with 'TAG=$(TAG)'
	git pull --tags
	git tag $(TAG)
	docker pull jmervine/noop-server:latest
	docker tag jmervine/noop-server:latest jmervine/noop-server:$(TAG)

.PHONY: require_tag
require_tag:
	@if test -z "$(TAG)"; then echo "TAG is required"; exit 1; fi

.PHONY: release
release: require_tag build push tag
	git push --tags
	docker push  jmervine/noop-server:$(TAG)

.PHONY: clean
clean:
	rm -rf bin

.PHONY: todos
todos:
	@git grep -n TODO | grep -v Makefile | awk -F':' '{ print " - TODO["$$1":"$$2"]:"$$NF }'