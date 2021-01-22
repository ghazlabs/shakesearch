.PHONY: build

build:
	docker build -t shakesearch -f ./build/package/search/Dockerfile .
run:
	make build
	docker run --rm \
		-v ${PWD}/cmd/search/web:/web \
		-p 3001:3001 \
		shakesearch