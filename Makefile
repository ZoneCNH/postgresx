.PHONY: fmt vet test race lint secret-scan security boundary contracts ci integration release-check

GO ?= go
GOENV ?= GOWORK=off

fmt:
	gofmt -w $$(find . -name '*.go' -not -path './.git/*')

vet:
	$(GOENV) $(GO) vet ./...

test:
	$(GOENV) $(GO) test ./...

race:
	$(GOENV) $(GO) test -race ./...

lint:
	bash ./scripts/ci/lint.sh

secret-scan:
	bash ./scripts/ci/secret_scan.sh

security: secret-scan

boundary:
	bash ./scripts/check_boundary.sh

contracts:
	bash ./scripts/check_contracts.sh

ci: fmt vet test race boundary contracts secret-scan

integration:
	bash ./scripts/ci/migration_up_down_up.sh

release-check:
	bash ./scripts/ci/release_check.sh
