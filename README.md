![Cover](./docs/assets/cover.png)

> **W**att **TF** is a powerful CLI tool that transforms JSON and YAML configurations into Terraform JSON files - Because writing complex nested loops in HCL shouldn't make you question your career choices


## Why?

Hi, I'm Sebastian :wave:

Like many engineers, I have a love-hate relationship with Terraform ( hate especially when dealing with complex HCL structures ). And I've also learned that most developers don't want to write Terraform just to deploy an application. 

The challenge isn't building infrastructure. It's making infrastructure easy to consume.

### Infrastructure with Dev Teams
A common setup looks like this:

The infrastructure team builds reusable Terraform modules.

The platform team defines standards, networking, security, IAM, monitoring, and compliance.

Development teams simply want to deploy their applications.

In theory, Terraform modules should make this easy.

In practice, developers still need to understand Terraform, module inputs, provider-specific options, and dozens of configuration parameters that they often shouldn't have to care about.

As organizations grow, every team starts building its own glue code.


### The Pain

Platform teams want consistency.

Developers want simplicity.

These goals often conflict.

Infrastructure teams want to expose a small, opinionated interface:

```json
{
  "service": "orders",
  "image": "ghcr.io/company/orders:v1",
  "database": true
}
```

Instead, developers end up interacting directly with complex Terraform modules containing dozens of variables, many of which should remain implementation details.

The result is duplicated tooling, inconsistent deployments, and platform teams spending valuable time maintaining custom generators instead of improving their infrastructure.

### The Solution: Watt TF

Watt TF is a generic transformation engine that sits between your application configuration and Terraform.

Instead of asking developers to write Terraform, platform teams define how structured input maps to reusable Terraform modules and resources.

With Watt TF, platform teams can:

- Build simple, stable interfaces for developers.
- Reuse existing Terraform modules without exposing unnecessary complexity.
- Apply organization-wide defaults and standards.
- Compose infrastructure from reusable building blocks.
- Generate consistent Terraform JSON from structured data.

Developers don't need to understand your Terraform implementation.

They only need to describe what they want.

Watt TF takes care of how it's built.

---

## Installation

### Via Go
```bash
go install github.com/devsebastianops/watt-tf/cmd/wtf@latest
wtf --help
```

### From Source
```bash
git clone https://github.com/devsebastianops/watt-tf.git
cd watt-tf
go build -o wtf ./cmd/wtf/main.go
./wtf --help
```

### Docker
```bash
docker run --rm -v "$(pwd):/app" -w /app devsebastianops/watt-tf --help
```

### From pre-built binaries
You can download the latest release from the [GitHub releases page](https://github.com/devsebastianops/watt-tf/releases).

---

## Quick Start

### Basic Usage

```bash
wtf build --input config.json --config .wtf.yaml --output terraform.tf.json
```

### Example 1: Simple JSON to Terraform

**input.json**
```json
{
  "cloudrun": {
    "image": "gcr.io/cloudrun/hello",
    "port": 8080
  }
}
```

**.wtf.yaml**
```yaml
transform:
  - target: resource.google_cloud_run_service.default
    value:
      image: "${input.cloudrun.image}"
      port: "${input.cloudrun.port}"
```

**Output (terraform.tf.json)**
```json
{
  "resource": {
    "google_cloud_run_service": {
      "default": {
        "image": "gcr.io/cloudrun/hello",
        "port": 8080
      }
    }
  }
}
```

Please refer to the [examples](example/README.md) for more detailed use cases and configurations, otherwise the [documentation](docs/README.md) provides comprehensive guidance.

---

## Command Reference

### `wtf build`

Transforms your input configuration into Terraform JSON.

```bash
wtf build [OPTIONS]

Options:
  --input <FILE>      Input file (JSON or YAML) [required]
  --config <FILE>     Transformation config (.wtf.yaml) [default: .wtf.yaml]
  --output <FILE>     Output file [default: watt.tf.json]
  --debug             Enable debug logging
  --strict            Enable strict mode (fail on missing keys)
  --schema <FILE>     JSON Schema file for input validation
```

**Examples:**
```bash
# Simple build
wtf build --input config.json

# Custom output
wtf build --input input.yaml --output main.tf.json

# Different config
wtf build --input config.json --config transformations.yaml --output output.tf.json

# Schema validation for input
wtf build --input input.json --config .wtf.yaml --schema schema.json --output output.tf.json
```

---

## Testing

Watt TF includes comprehensive E2E tests covering all features:

```bash
# Run all tests
go test ./tests/e2e -v

# Run specific example
go test ./tests/e2e -v -run "TestE2EExamples/01-simple-json"

# View test examples
ls -la example/
```

--- 

### CI/CD Integration
Generate Terraform configurations dynamically in your pipeline:

```bash
#!/bin/bash
# Extract environment from Git branch
ENV=$(git rev-parse --abbrev-ref HEAD)

# Generate configuration
wtf build \
  --input environment.json \
  --config transforms.yaml \
  --output terraform-${ENV}.tf.json

# Apply Terraform
terraform init
terraform plan -out=plan.tfplan
terraform apply plan.tfplan
```

---

## Feedback & Contribution

We ❤️ contributions! Whether it's bug reports, feature requests, or pull requests:

### Report Issues
Found a bug? [Open an Issue](https://github.com/devsebastianops/watt-tf/issues)

### Contribute Code
1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Commit changes: `git commit -m 'Add amazing feature'`
4. Push to branch: `git push origin feature/amazing-feature`
5. Open a Pull Request


---

## License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

---

## Contact & Support

- **GitHub**: [devsebastianops/watt-tf](https://github.com/devsebastianops/watt-tf)
- **Issues**: [GitHub Issues](https://github.com/devsebastianops/watt-tf/issues)

---

Made with ❤️, **Happy Terraforming!** 🚀
