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

.PHONY: clean
clean:
	rm -rf bin