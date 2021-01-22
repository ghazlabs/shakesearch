.PHONY: build

build:
	docker build -t shakesearch -f ./build/package/search/Dockerfile .
run:
	make build
	docker run --rm -p 3001:3001 shakesearch