all: build

build:
	docker build . --target build \
	--platform local \
	--output ./bin

e2e:
	docker build . --target e2e

.PHONY: build e2e
