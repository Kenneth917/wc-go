IMAGE ?= wc-go

.PHONY: test unit bdd image e2e
test: unit bdd

unit:
	go test -race -cover ./internal/...

bdd:
	go test ./features/...

image:
	docker build -t $(IMAGE) .

# Same scenarios, run against the container instead of a local binary.
e2e: image
	WC_DOCKER_IMAGE=$(IMAGE) go test ./features/...
