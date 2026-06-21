.PHONY: build run run-local test dashboards alerts docs

ROOT := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
ARGS ?=

build:
	docker build -t fantastic-beasts "$(ROOT)"

run:
	docker run --rm -p 8080:8080 --name fantastic-beasts fantastic-beasts

run-local:
	cd "$(ROOT)src" && go run .

test:
	cd "$(ROOT)src" && go build -o /dev/null .

dashboards:
	"$(ROOT)cmd/apply-dashboards" $(ARGS)

alerts:
	"$(ROOT)cmd/apply-alerts" $(ARGS)

docs:
	@echo "Docs: $(ROOT)docs/README.md"
	@echo "Ops:  $(ROOT)docs/ops/README.md"
