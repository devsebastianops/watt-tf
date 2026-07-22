# Quick start

This guide will help you getting started with Watt TF quickly. 

To make you understand the basics, we will create a simple terraform resource that aims to create a local file with some content, nothing terraform itself would not be able to do by itself. This is just a simple example to getting you started, but it will give you a good understanding of how Watt TF works.

For more complex use cases, please refer to the [real world scenarios](../examples/overview.md).

## 1. Install Watt TF

::: code-group
```sh [Installer]
curl -sSL https://raw.githubusercontent.com/devsebastianops/watt-tf/main/install.sh | bash
```

```sh [Go]
go install github.com/devsebastianops/watt-tf/cmd/wtf@latest
```

```sh [From Source]
git clone https://github.com/devsebastianops/watt-tf.git
cd watt-tf
go build -o wtf ./cmd/wtf/main.go
```
:::

<br>

::: tip
We recommend to use the installer script, as it will download the latest release and put it in your path.
:::

## 2. Create your input data

::: code-group
```json [JSON]
{
    "name": "microservice",
    "content": "This is a microservice"
}
```

```yaml [YAML]
name: microservice
content: |
  This is a microservice
```
:::

## 3. Create your blueprint

```yaml
transform:
- target: resource.local_file.${input.name}
  value:
    filename: ${input.name}
    content: ${input.content}
```

<br>

::: tip
We also provide a [schema](/schema/watt-tf-configuration.schema.json) for the blueprint configuration, which can be used for validation and autocompletion in your IDE.
:::

## 4. Run the transformation

```sh
wtf build -i input.json -b blueprint.yaml -o output.tf.json
```

The `output.tf.json` should look like this:

```json
{
  "resource": {
    "local_file": {
      "microservice": {
        "filename": "microservice",
        "content": "This is a microservice"
      }
    }
  }
}
```

This is valid terraform json ( `.tf.json` ) that can be applied.

## 5. Let terraform handle what it does best

```sh
terraform init
terraform plan -out planfile
terraform apply planfile
```