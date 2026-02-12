# Makefile for vela-go-definitions

# Go parameters
GOCMD=go
GOMOD=$(GOCMD) mod

# Ginkgo parameters
GINKGO=$(shell which ginkgo 2>/dev/null || echo "go run github.com/onsi/ginkgo/v2/ginkgo")

# Test data path
TESTDATA_PATH ?= test/builtin-definition-example

# Vela CLI path (can be overridden)
VELA_CLI ?= vela

# Timeout for E2E tests
E2E_TIMEOUT ?= 30m

.PHONY: tidy install-ginkgo test-e2e test-e2e-components test-e2e-traits test-e2e-policies test-e2e-workflowsteps help

## Dependency management
tidy:
	@echo "Tidying dependencies..."
	$(GOMOD) tidy

## Install Ginkgo CLI
install-ginkgo:
	@echo "Installing Ginkgo CLI..."
	go install github.com/onsi/ginkgo/v2/ginkgo@latest

## E2E Test targets
test-e2e: test-e2e-components test-e2e-traits test-e2e-policies test-e2e-workflowsteps
	@echo "All E2E tests completed!"

test-e2e-components:
	@echo "Running E2E tests for component definitions..."
	TESTDATA_PATH=$(TESTDATA_PATH) VELA_CLI=$(VELA_CLI) \
		$(GINKGO) -v --timeout=$(E2E_TIMEOUT) --label-filter="components" ./test/e2e/...

test-e2e-traits:
	@echo "Running E2E tests for trait definitions..."
	TESTDATA_PATH=$(TESTDATA_PATH) VELA_CLI=$(VELA_CLI) \
		$(GINKGO) -v --timeout=$(E2E_TIMEOUT) --label-filter="traits" ./test/e2e/...

test-e2e-policies:
	@echo "Running E2E tests for policy definitions..."
	TESTDATA_PATH=$(TESTDATA_PATH) VELA_CLI=$(VELA_CLI) \
		$(GINKGO) -v --timeout=$(E2E_TIMEOUT) --label-filter="policies" ./test/e2e/...

test-e2e-workflowsteps:
	@echo "Running E2E tests for workflowstep definitions..."
	TESTDATA_PATH=$(TESTDATA_PATH) VELA_CLI=$(VELA_CLI) \
		$(GINKGO) -v --timeout=$(E2E_TIMEOUT) --label-filter="workflowsteps" ./test/e2e/...

## Help
help:
	@echo "Available targets:"
	@echo "  tidy                   - Tidy go.mod dependencies"
	@echo "  install-ginkgo         - Install Ginkgo CLI for running E2E tests"
	@echo "  test-e2e               - Run all E2E tests"
	@echo "  test-e2e-components    - Run E2E tests for component definitions"
	@echo "  test-e2e-traits        - Run E2E tests for trait definitions"
	@echo "  test-e2e-policies      - Run E2E tests for policy definitions"
	@echo "  test-e2e-workflowsteps - Run E2E tests for workflowstep definitions"
	@echo ""
	@echo "Environment variables:"
	@echo "  TESTDATA_PATH - Path to test data (default: test/builtin-definition-example)"
	@echo "  VELA_CLI      - Path to vela CLI (default: vela)"
	@echo "  E2E_TIMEOUT   - Timeout for E2E tests (default: 30m)"

