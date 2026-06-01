.PHONY: fmt vet test race lint secret-scan security boundary contracts ci ci-extended integration evidence release-evidence-check release-preflight release-final-check release-check

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

ci-extended: ci integration release-evidence-check

integration:
	bash ./scripts/ci/migration_up_down_up.sh

evidence:
	bash ./scripts/ci/evidence.sh

release-evidence-check:
	bash ./scripts/ci/release_evidence_check.sh

release-preflight: ci-extended evidence

release-final-check: release-preflight
	git diff --check

release-check:
	bash ./scripts/ci/release_check.sh
