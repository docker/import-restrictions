build:
	go build -o /out/import-restrictions .

e2e:
	go test ./e2e

.PHONY: build e2e
