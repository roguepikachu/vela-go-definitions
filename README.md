# vela-go-definitions

A collection of KubeVela X-Definitions written in Go using the [defkit](https://github.com/oam-dev/kubevela/tree/master/pkg/definition/defkit) fluent builder API.

## Overview

This module contains Go-based KubeVela X-Definitions that can be applied to any KubeVela cluster. Definitions are written in Go and generate CUE automatically via the defkit framework.

## Directory Structure

- **components/** - ComponentDefinitions for workload types
- **traits/** - TraitDefinitions for operational behaviors
- **policies/** - PolicyDefinitions for application policies
- **workflowsteps/** - WorkflowStepDefinitions for delivery workflows
- **cmd/defkit/** - CLI tool for generating CUE and exporting definitions
- **vela-templates/definitions/** - Generated CUE output (do not edit manually)

## Usage

### Apply all definitions

```bash
vela def apply-module github.com/oam-dev/vela-go-definitions
```

### List definitions

```bash
vela def list-module github.com/oam-dev/vela-go-definitions
```

### Validate definitions

```bash
vela def validate-module github.com/oam-dev/vela-go-definitions
```

### Apply with namespace

```bash
vela def apply-module github.com/oam-dev/vela-go-definitions --namespace my-namespace
```

### Dry-run (preview without applying)

```bash
vela def apply-module github.com/oam-dev/vela-go-definitions --dry-run
```

## Development

### Prerequisites

- Go 1.23+
- [golangci-lint](https://golangci-lint.run/) (installed automatically by `make lint` if missing)

### Make Reviewable

Before submitting a pull request, run `make reviewable` to ensure your changes are ready for review. This is the same pattern used in the [kubevela](https://github.com/oam-dev/kubevela) repository.

```bash
make reviewable
```

This runs the following checks in order:

1. **generate** - Regenerates CUE definitions from Go into `vela-templates/definitions/`
2. **fmt** - Formats all Go code
3. **vet** - Runs `go vet` on all packages
4. **lint** - Runs `golangci-lint`
5. **check-diff** - Verifies no uncommitted changes in generated files

If `check-diff` fails, it means the generated CUE files are out of date. Run `make generate` and commit the updated files.

### Individual Targets

```bash
make generate    # Regenerate CUE definitions
make fmt         # Format Go code
make vet         # Vet Go code
make lint        # Lint Go code
make check-diff  # Verify generated files are up-to-date
make tidy        # Tidy go.mod dependencies
```

### CLI Tool

The `cmd/defkit` CLI provides two subcommands:

```bash
# Generate CUE files into vela-templates/definitions/
go run ./cmd/defkit generate

# Export all registered definitions as JSON
go run ./cmd/defkit register
```

## Adding New Definitions

1. Create a new Go file in the appropriate directory
2. Add an `init()` function that registers your definition
3. Use the defkit package fluent API to define your component/trait/policy/workflow-step
4. Run `make reviewable` to regenerate CUE and validate

Example component definition:

```go
package components

import "github.com/oam-dev/kubevela/pkg/definition/defkit"

func init() {
    defkit.Register(MyComponent())
}

func MyComponent() *defkit.ComponentDefinition {
    image := defkit.String("image").Mandatory().Description("Container image")
    replicas := defkit.Int("replicas").Default(1).Description("Number of replicas")

    return defkit.NewComponent("my-component").
        Description("My custom component").
        Workload("apps/v1", "Deployment").
        Params(image, replicas).
        Template(myComponentTemplate)
}

func myComponentTemplate(tpl *defkit.Template) {
    vela := defkit.VelaCtx()
    image := defkit.String("image")
    replicas := defkit.Int("replicas")

    deployment := defkit.NewResource("apps/v1", "Deployment").
        Set("spec.replicas", replicas).
        Set("spec.selector.matchLabels[app.oam.dev/component]", vela.Name()).
        Set("spec.template.spec.containers[0].name", vela.Name()).
        Set("spec.template.spec.containers[0].image", image)

    tpl.Output(deployment)
}
```

## Testing

### Unit Tests

```bash
make test-unit
```

### E2E Tests

E2E tests validate definitions against a live KubeVela cluster.

```bash
# Install Ginkgo CLI
make install-ginkgo

# Run all E2E tests
make test-e2e

# Run specific test categories
make test-e2e-components    # Component definitions
make test-e2e-traits        # Trait definitions
make test-e2e-policies      # Policy definitions
make test-e2e-workflowsteps # Workflow step definitions
```

### Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PROCS` | 10 | Number of parallel test processes |
| `E2E_TIMEOUT` | 10m | Test timeout duration |
| `TESTDATA_PATH` | `test/builtin-definition-example` | Path to test data |
| `DEFINITIONS_DIR` | `vela-templates/definitions` | Output directory for generated CUE |

## CI/CD

GitHub Actions workflows automatically run on push and pull requests to `main`:

- **Unit Tests** - Runs all unit tests
- **Reviewable** - Runs `make generate`, `make fmt`, `make vet`, linting, and `check-diff` to ensure generated files are up-to-date

## License

Apache License 2.0
