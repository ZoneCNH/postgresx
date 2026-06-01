.PHONY: fmt vet test race lint govulncheck secret-scan security boundary contracts ci ci-extended integration evidence release-check release-preflight release-evidence-check release-final-check

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
	bash ./scripts/ci/secret_scan.sh

security: secret-scan govulncheck

boundary:
	bash ./scripts/check_boundary.sh

contracts:
	bash ./scripts/check_contracts.sh

ci: fmt vet test race boundary contracts secret-scan lint

ci-extended: ci integration release-evidence-check

ci-extended: ci integration release-evidence-check

integration:
	bash ./scripts/ci/migration_up_down_up.sh

evidence:
	bash ./scripts/ci/write_evidence.sh "$(VERSION)"

release-check:
	bash ./scripts/ci/release_check.sh "$(VERSION)"

release-preflight: ci-extended release-check

release-evidence-check:
	bash ./scripts/ci/release_evidence_check.sh "$(VERSION)"

release-final-check: release-preflight
	git diff --check
	@test -z "$$(git status --short)" || (echo "release-final-check requires a clean worktree" >&2; git status --short >&2; exit 1)
