<div align="center">
  <img src="./docs/public/assets/watt-tf-mascott-sticker-outlines.png" alt="Watt TF Mascot" width="180" />
  <h1>Watt TF</h1>
  <p><strong>Terraform without the HCL gymnastics.</strong></p>
  <p>
    <a href="https://github.com/devsebastianops/watt-tf/actions"><img src="https://img.shields.io/github/actions/workflow/status/devsebastianops/watt-tf/on-pull-request-ci.yaml?branch=main&style=flat-square" alt="CI Status"></a>
    <a href="https://github.com/devsebastianops/watt-tf/releases"><img src="https://img.shields.io/github/v/release/devsebastianops/watt-tf?style=flat-square" alt="Latest Release"></a>
    <a href="https://github.com/devsebastianops/watt-tf/blob/main/LICENSE"><img src="https://img.shields.io/github/license/devsebastianops/watt-tf?style=flat-square" alt="License"></a>
  </p>
  <h3>
    <a href="https://watt-tf.dev">Documentation</a>
    <span> | </span>
    <a href="https://watt-tf.dev/guide/quick-start">Quick Start</a>
    <span> | </span>
    <a href="https://watt-tf.dev/examples/overview">Real-World Examples</a>
  </h3>
</div>

---

**Watt TF** is a declarative blueprint engine that compiles raw data (JSON/YAML) into clean, standard Terraform JSON. 

Instead of forcing your developers to fight complex HCL syntax, or asking your platform team to maintain brittle custom templating scripts, **Platform Engineers** define reusable blueprints using Google CEL (Common Expression Language), while **Developers** just supply simple data.

---

## Why Watt TF? ⚡

- **No HCL Gymnastics:** Stop writing unreadable nested `for`/`for_each` loops, `dynamic` blocks, and ternary string-splitting in Terraform.
- **100% Type-Safe JSON:** No text-replacement hacks (like Jinja2 or Go templates) that generate broken syntax. We compile directly to Terraform JSON.
- **Platform-Engineering Native:** Build clean, opinionated interfaces for your developers to deploy apps without exposing them to the raw horror of naked infrastructure code.

---

## Quick Look: How it Works

### 1. Simple Data ( JSON or YAML )

Developers supply only what they care about:

```json
{
  "service": {
    "name": "ms-orders",
  },
  "app": {
    "image": "ghcr.io/acme/ms-orders",
    "version": "v1.0.0",
    "port": 8080
  },
  "database": {
    "type": "postgres",
    "version": "14",
    "size": "20Gi"
  }
}
```

### 2. The Platform Blueprint

The platform team defines how that data translates to Terraform modules, using powerful embedded Google CEL logic:

```yaml
transform:
  - target: module.app.${input.service.name}
    value: 
        source: "git::https://git.acme.com/terraform-modules/app.git?ref=v1.0.0"
        version: "${input.app.version}"
        port: "${input.app.port}"
        min_replicas: ${has(input.app.min_replicas) ? optional(input.app.min_replicas) : optional.none()}
        max_replicas: ${has(input.app.max_replicas) ? optional(input.app.max_replicas) : optional.none()}
        ram: ${env.ENVIRONMENT == "production" ? "2Gi" : "512Mi"}
        database: ${has(input.database) ? "module.database.${input.service.name}.dsn" : optional.none()}

  - target: module.database.${input.service.name}
    if: has(input.database)
    value:
        source: "git::https://git.acme.com/terraform-modules/database.git?ref=v1.0.0"
        type: "${input.database.type}"
        version: "${input.database.version}"
        size: "${input.database.size}"
```

### 3. Transform to Terraform JSON

```bash
wtf build --input input.json --config blueprint.yaml --output terraform.json
```

#### Expected Result

```json
{
  "module": {
    "app": {
      "ms-orders": {
        "source": "git::https://git.acme.com/terraform-modules/app.git?ref=v1.0.0",
        "version": "v1.0.0",
        "port": 8080,
        "min_replicas": null,
        "max_replicas": null,
        "ram": "512Mi",
        "database": "module.database.ms-orders.dsn"
      }
    },
    "database": {
      "ms-orders": {
        "source": "git::https://git.acme.com/terraform-modules/database.git?ref=v1.0.0",
        "type": "postgres",
        "version": "14",
        "size": "20Gi"
      }
    }
  }
}
```

---

## Getting Started

The fastest way to install and try out Watt TF locally:

```bash
curl -sSL https://raw.githubusercontent.com/devsebastianops/watt-tf/main/install.sh | bash
wtf --help
```

For more installation methods and the full documentation check out [watt-tf.dev](https://watt-tf.dev).

---

## Examples

All configurations in the `/examples` directory are actively tested in our E2E pipeline on every commit.

You can use these to understand how Watt TF works, or as a starting point for your own blueprints.

> [!TIP]
> Have a look into our real world examples at [watt-tf.dev/examples/overview](https://watt-tf.dev/examples/overview).

---

## Community, Feedback & Support

We ❤️ contributions! Whether you found a bug, want to propose a new feature, or improve the codebase:

*   **Report Bugs:** Found an issue? [Open an Issue on GitHub](https://github.com/devsebastianops/watt-tf/issues).
*   **Contribute:** Check out our [Contribution Guidelines](CONTRIBUTING.md) to see how to get started with the codebase and documentation.
*   **Security:** If you discovered a security vulnerability, please refer to our [Security Policy](SECURITY.md) for responsible disclosure instructions.

---

## License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.


---

<p align="center">
  Made with ❤️ by <a href="https://github.com/devsebastianops">devsebastianops</a>. Happy Terraforming! 🚀
</p>
