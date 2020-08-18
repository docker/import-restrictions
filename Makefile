all: build

build:
	docker build . --target build \
	--platform local \
	--output ./bin

.PHONY: build
