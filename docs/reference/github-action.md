# Watt TF Build GitHub Action

Watt TF provides a GitHub Action to streamline the process of building your Watt TF blueprints directly from your GitHub repository. 

You can find the GitHub Action in the [GitHub Marketplace](https://github.com/marketplace/actions/watt-tf-action) or in its [official repository](https://github.com/devsebastianops/watt-tf-build-action).

## Inputs

| Name | Description | Required | Default |
|------|-------------|----------|---------|
| `version` | The version of Watt TF to use (e.g. latest or v1.0.0) | `false` | `'latest'` |
| `input` | The input file as JSON or YAML | `true` | |
| `blueprint` | The blueprint with the transformation rules | `true` | |
| `output` | Where should Watt TF put the .tf.json result | `true` | |
| `schema` | Optional JSON schema file for input validation | `false` | |
| `strict` | Run Watt TF build in strict mode | `false` | `'false'` |
| `stripNulls` | Run Watt TF build with strip-nulls flag | `false` | `'false'` |


## Example Usage

```yaml
  uses: devsebastianops/watt-tf-build-action@v1
  with:
    version: 'latest'
    input: '/path/to/input.json'
    blueprint: '/path/to/blueprint.yaml'
    output: '/path/to/output.tf.json'
    strict: 'true'
    stripNulls: 'true'
```

