.PHONY: fmt vet test race lint secret-scan security boundary contracts foundationx-api template-alignment ci ci-extended integration evidence release-check release-preflight release-evidence-check release-final-check

GO ?= go
GOENV ?= GOWORK=off
VERSION ?= v0.1.0

fmt:
	gofmt -w $$(find . -name '*.go' -not -path './.git/*')

vet:
	$(GOENV) $(GO) vet ./...

test:
	$(GOENV) $(GO) test ./...

race:
	$(GOENV) $(GO) test -race ./...

govulncheck:
	bash ./scripts/ci/govulncheck.sh

lint:
	bash ./scripts/ci/lint.sh

secret-scan:
	bash ./scripts/check_secrets.sh

security: secret-scan govulncheck

boundary:
	bash ./scripts/check_boundary.sh

contracts:
	bash ./scripts/check_contracts.sh

foundationx-api:
	bash ./scripts/check_foundationx_api.sh

template-alignment:
	bash ./scripts/check_template_alignment.sh

ci: fmt vet test race boundary contracts secret-scan

ci-extended: ci lint security foundationx-api template-alignment

integration:
	bash ./scripts/run_integration.sh

evidence:
	bash ./scripts/generate_manifest.sh

release-check:
	bash ./scripts/ci/release_check.sh

release-preflight: ci-extended integration evidence

release-evidence-check:
	bash ./scripts/ci/release_evidence_check.sh $(VERSION)

release-final-check: ci-extended release-evidence-check
	$(GOENV) $(GO) list -m | grep -Fx github.com/ZoneCNH/postgresx
	$(GOENV) $(GO) list ./pkg/postgresx >/dev/null
	git diff --check
