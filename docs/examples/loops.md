# Loops

This section demonstrates how to define and use loops within Watt TF blueprints. Loops allow you to iterate over lists or maps and dynamically generate configuration blocks based on the items in the collection.

## Scenario

We want to create multiple terraform resources for storage buckets, based on a list of bucket definitions provided as input.

## Solution

### Input

```yaml
buckets:
  - name: order-export
    region: eu-west-1
  - name: product-export
    region: eu-west-2
```

### Blueprint

Note, that `item_index` is available within loops and can be used to uniquely identify each iteration by its zero based index. The `item` represents the current element in the collection being iterated over.

```yaml
transform:
  - target: module.bucket-${item_index}
    for_each: "${input.buckets}"
    value:
        source: "git::https://github.com/myorg/myrepo.git?ref=v1.2.3"
        name: "${item.name}"
        region: "${item.region}"
```

### Command

```bash
wtf build --input input.yaml --blueprint blueprint.yaml --output build.tf.json
```

### Output

```json
{
    "module": {
        "bucket-0": {
            "source": "git::https://github.com/myorg/myrepo.git?ref=v1.2.3",
            "name": "order-export",
            "region": "eu-west-1"
        },
        "bucket-1": {
            "source": "git::https://github.com/myorg/myrepo.git?ref=v1.2.3",
            "name": "product-export",
            "region": "eu-west-2"
        }
    }
}
```
