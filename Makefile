XLIB_CONTEXT ?= local_write
GO ?= go
GOENV ?= GOWORK=off
VERSION ?= v1.0.0

.PHONY: require-gowork-off
require-gowork-off:
	@if [ "$${GOWORK:-}" != "off" ]; then \
		echo "GOWORK=off is required for release targets"; \
		exit 1; \
	fi

.PHONY: build
build:
	$(GOENV) $(GO) build ./...

.PHONY: build-check
build-check: build

.PHONY: shell
shell:
	bash

.PHONY: fmt
fmt:
	gofmt -w $$(find . -name '*.go' -not -path './.git/*')

.PHONY: vet
vet:
	$(GOENV) $(GO) vet ./...

.PHONY: test
test:
	$(GOENV) $(GO) test ./...

.PHONY: test-unit
test-unit:
	$(GOENV) $(GO) test ./pkg/postgresx ./testkit ./contracts

.PHONY: test-contract
test-contract:
	$(GOENV) $(GO) test ./test/contract

.PHONY: test-integration
test-integration:
	POSTGRESX_REQUIRE_INTEGRATION=1 bash ./scripts/run_integration.sh

.PHONY: test-chaos
test-chaos:
	$(GOENV) $(GO) test ./test/chaos

.PHONY: benchmark-smoke
benchmark-smoke:
	bash ./scripts/ci/benchmark_smoke.sh

.PHONY: downstream-smoke
downstream-smoke:
	bash ./scripts/ci/downstream_smoke.sh

.PHONY: race
race:
	$(GOENV) $(GO) test -race ./...

.PHONY: lint
lint:
	bash ./scripts/ci/lint.sh

.PHONY: integration
integration: test-integration

.PHONY: integration-check
integration-check: integration

.PHONY: secret-scan
secret-scan:
	bash ./scripts/check_secrets.sh

.PHONY: govulncheck
govulncheck:
	bash ./scripts/ci/govulncheck.sh

.PHONY: security
security: secret-scan govulncheck

.PHONY: boundary
boundary:
	bash ./scripts/check_boundary.sh

.PHONY: boundary-check
boundary-check: boundary

.PHONY: contracts
contracts:
	bash ./scripts/check_contracts.sh

.PHONY: contract-check
contract-check: contracts

.PHONY: foundationx-api
foundationx-api:
	bash ./scripts/check_foundationx_api.sh

.PHONY: template-alignment
template-alignment:
	bash ./scripts/check_template_alignment.sh

.PHONY: docs-check
docs-check:
	bash ./scripts/ci/release_evidence_check.sh $(VERSION)

.PHONY: evidence
evidence:
	bash ./scripts/generate_manifest.sh
	python3 ./scripts/ci/l2_evidence.py --version $(VERSION)

.PHONY: l2-evidence
l2-evidence:
	python3 ./scripts/ci/l2_evidence.py --version $(VERSION)

.PHONY: release-check
release-check:
	bash ./scripts/ci/release_check.sh

.PHONY: release-preflight
release-preflight: ci-extended integration evidence

.PHONY: release-evidence-check
release-evidence-check:
	bash ./scripts/ci/release_evidence_check.sh $(VERSION)

.PHONY: release-final-check
release-final-check: release-evidence-check test
	$(GOENV) $(GO) list -m | grep -Fx github.com/ZoneCNH/postgresx
	$(GOENV) $(GO) list ./pkg/postgresx >/dev/null
	git diff --check
	@if [ -n "$$(git status --short)" ]; then \
		echo "release-final-check requires a clean workspace; run make evidence and commit artifacts first" >&2; \
		git status --short >&2; \
		exit 1; \
	fi

.PHONY: property
property:
	go test ./... -run 'Test.*Property|Test.*Invariant'

.PHONY: golden
golden:
	go test ./... -run 'Test.*Golden|Test.*Snapshot'

.PHONY: install-hooks
install-hooks:
	@git config core.hooksPath .githooks
	@echo "✅ git hooks 已启用（core.hooksPath=.githooks）"

.PHONY: doctor-hooks
doctor-hooks:
	@[ "$$(git config --get core.hooksPath)" = ".githooks" ] || { \
	  echo "ERROR: core.hooksPath 未指向 .githooks，请运行 make install-hooks"; \
	  exit 1; \
	}
	@echo "✅ hooks 配置正确"

.PHONY: doctor-hooks-local
doctor-hooks-local:
	@if [ -n "$$CI" ] || [ -n "$$GITHUB_ACTIONS" ]; then \
	  echo "doctor-hooks-local: CI 环境，跳过 hooks 检查"; \
	else \
	  [ "$$(git config --get core.hooksPath)" = ".githooks" ] || { \
	    echo "ERROR: 本地 git hooks 未启用 (core.hooksPath != .githooks)"; \
	    echo "       这会跳过 pre-commit secret 扫描，请运行: make install-hooks"; \
	    exit 1; \
	  }; \
	  echo "✅ 本地 hooks 已启用"; \
	fi

.PHONY: sync-main
sync-main:
	@git fetch origin main
	@CUR=$$(git rev-parse --abbrev-ref HEAD); \
	if [ "$$CUR" = "main" ]; then \
	  git merge --ff-only origin/main && echo "✅ main 已同步至 origin/main"; \
	else \
	  LOCAL=$$(git rev-parse main 2>/dev/null || echo ""); \
	  REMOTE=$$(git rev-parse origin/main); \
	  if [ "$$LOCAL" = "$$REMOTE" ]; then \
	    echo "✅ 本地 main 已是最新（当前在 $$CUR）"; \
	  else \
	    echo "⚠️  当前不在 main 分支（$$CUR）"; \
	    echo "   本地 main: $$LOCAL"; \
	    echo "   远端 main: $$REMOTE"; \
	    echo "   请到主 worktree 执行: git merge --ff-only origin/main"; \
	    exit 1; \
	  fi; \
	fi

# ── Docker Toolchain ──────────────────────────────────────────
DOCKER_IMAGE ?= $(notdir $(CURDIR))-toolchain:local
DOCKER_GATE ?= ./scripts/docker/docker_gate.sh

.PHONY: docker-toolchain-check
docker-toolchain-check:
	./scripts/docker/check_toolchain.sh

.PHONY: docker-build
docker-build: docker-toolchain-check
	DOCKER_BUILDKIT=1 docker buildx build --load --target toolchain --build-arg GO_VERSION=$${GO_VERSION:-1.25} --build-arg GOLANGCI_LINT_VERSION=$${GOLANGCI_LINT_VERSION:-v2.1.6} --build-arg GOVULNCHECK_VERSION=$${GOVULNCHECK_VERSION:-v1.1.4} --tag $(DOCKER_IMAGE) .

.PHONY: docker-shell
docker-shell: docker-build
	docker run --rm -it \
		--add-host host.docker.internal:host-gateway \
		--workdir /workspace \
		--volume "$(CURDIR):/workspace" \
		--volume go-build-cache:/root/.cache/go-build \
		--volume go-mod-cache:/go/pkg/mod \
		--env "GOWORK=$${GOWORK:-off}" \
		--env "XLIB_CONTEXT=$${XLIB_CONTEXT:-docker_toolchain}" \
		--env "VERSION=$${VERSION:-}" \
		--env "DOWNSTREAM=$${DOWNSTREAM:-}" \
		--env "XLIB_ENABLE_VULNCHECK=$${XLIB_ENABLE_VULNCHECK:-}" \
		--env "XLIB_FORCE_VULNCHECK=$${XLIB_FORCE_VULNCHECK:-}" \
		--env "XLIB_VULNCHECK_INTERVAL_HOURS=$${XLIB_VULNCHECK_INTERVAL_HOURS:-}" \
		--env "POSTGRES_TEST_DSN=$${POSTGRES_TEST_DSN:-}" \
		--env "POSTGRESX_INTEGRATION_DSN=$${POSTGRESX_INTEGRATION_DSN:-}" \
		--env "POSTGRESX_REQUIRE_INTEGRATION=$${POSTGRESX_REQUIRE_INTEGRATION:-}" \
		--env "CI=$${CI:-}" \
		--env "GITHUB_ACTIONS=$${GITHUB_ACTIONS:-}" \
		--env "GIT_CONFIG_COUNT=1" \
		--env "GIT_CONFIG_KEY_0=safe.directory" \
		--env "GIT_CONFIG_VALUE_0=/workspace" \
		$(DOCKER_IMAGE) bash

.PHONY: docker-ci
docker-ci:
	$(DOCKER_GATE) ci

.PHONY: docker-build-check
docker-build-check:
	$(DOCKER_GATE) build-check

.PHONY: docker-release-check
docker-release-check:
	$(DOCKER_GATE) release-check

.PHONY: docker-release-final-check
docker-release-final-check:
	$(DOCKER_GATE) release-final-check

.PHONY: runtime-check
runtime-check: build-check

.PHONY: docker-runtime-check
docker-runtime-check:
	$(DOCKER_GATE) runtime-check

.PHONY: drift-check
drift-check: template-alignment

.PHONY: docker-drift-check
docker-drift-check:
	$(DOCKER_GATE) drift-check

.PHONY: docker-contract
docker-contract: docker-toolchain-check docker-build-check docker-runtime-check docker-drift-check

# ── Composite Targets ─────────────────────────────────────────
.PHONY: ci
ci: doctor-hooks-local fmt vet test race boundary contracts secret-scan

.PHONY: ci-extended
ci-extended: ci foundationx-api template-alignment
