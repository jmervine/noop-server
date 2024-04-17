default: test build

.PHONY: test
test:
	go test -race ./...

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
	go run ./cmd/noop-server/ --record

bin/noop-server:
	go build -o bin/noop-server ./cmd/noop-server/...

# Should run 'make build push' first.
.PHONY: tag
tag: require_owner require_tag
	# Tag git and docker with 'TAG=$(TAG)'
	git pull --tags
	git tag $(TAG)
	docker pull jmervine/noop-server:latest
	docker tag jmervine/noop-server:latest jmervine/noop-server:$(TAG)

.PHONY: require_tag
require_tag:
	@if test -z "$(TAG)"; then echo "TAG is required"; exit 1; fi

# This is ugly and fragile, but meh, it'll work for now.
.PHONY: require_owner
require_owner:
	@# This is ugly and fragile, but meh, it'll work for now.
	@if ! [ "$(USER)" = "jmervine" ]; then echo "You are not the owner and cannot perform this task."; exit 1; fi

.PHONY: release
release: require_owner require_tag build push tag
	git push --tags
	docker push  jmervine/noop-server:$(TAG)

.PHONY: clean
clean:
	rm -rf bin

.PHONY: todos
todos:
	@git grep -n TODO | grep -v Makefile | awk -F':' '{ print " - TODO["$$1":"$$2"]:"$$NF }'
