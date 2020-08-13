CMD=?

lint:
	./scripts/test/lint
.PHONY: lint

test:
	./scripts/test/test
.PHONY: test

build:
	./scripts/build
.PHONY: build

terraform:
	./scripts/test/terraform $(CMD)
.PHONY: terraform
