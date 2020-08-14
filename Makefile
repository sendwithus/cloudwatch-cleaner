CMD=?

lint:
	./scripts/test/lint
.PHONY: lint

unit:
	./scripts/test/unit
.PHONY: unit

build:
	./scripts/build
.PHONY: build

terraform:
	./scripts/test/terraform $(CMD)
.PHONY: terraform
