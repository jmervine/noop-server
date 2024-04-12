default: test build

test:
	go test -v -race ./...

build:
	docker build -t jmervine/noop-server:latest .
	docker tag jmervine/noop-server:latest jmervine/noop-server:$(shell git reflog | head -n1 | cut -d' ' -f1)

push:
	docker push jmervine/noop-server:latest
	docker push jmervine/noop-server:$(shell git reflog | head -n1 | cut -d' ' -f1)

run:
	go run .