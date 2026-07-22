# Includes

Especially when dealing with large blueprints, it is often useful to split your blueprint into multiple files. Watt TF supports this by allowing you to include other files in your blueprint.

## Including other blueprints

Your blueprint can include other blueprints using the `include` keyword. The included file will be processed as if it were part of the main blueprint.

```yaml
include:
  - compute.yaml
```

Your blueprint can not just include other blueprints, but can contain regular transform blocks alongside the include statement. The included blueprint will be processed in the order it appears in the transform array.

```yaml
include:
  - compute.yaml

transform:
  - target: resource.google_cloud_run_service.default
    value:
      port: ${input.cloudrun.port}
      image: ${input.cloudrun.image}

```

## Recursive Includes

Currently, Watt TF does not support recursive includes. If an included blueprint contains an `include` statement that references another blueprint, it will be ignored.

There are - at least for now -multiple reasons for this limitation. 

### 1. Explicitness

```bash
blueprint.yaml
├── compute.yaml
├── network.yaml
└── iam.yaml
```

With the current implementation it is crystal clear that there can be only one main blueprint and all other blueprints are included explicitly. This makes it easier to reason about the blueprint structure and avoid accidental includes.


```bash
blueprint.yaml
└── project.yaml
    ├── globals.yaml
    │   ├── providers.yaml
    │   └── locals.yaml
    └── services.yaml
        ├── api.yaml
        └── workers.yaml
            └── monitoring.yaml
```

For the example above, it is not immediately clear which blueprints are included and which are not. This can lead to confusion and make it harder to reason about the blueprint structure - You will have to open multiple files to understand the blueprint structure.

That also leads to hidden dependencies, which can be a source of bugs and make it harder to maintain the blueprint over time.

### 2. The right order

Recursive includes can lead to questions about the order of the execution.

Given is the following example:

```bash
blueprint.yaml
├── module-a.yaml
└── module-b.yaml
```

```bash
module-a.yaml
├── module-c.yaml
└── module-d.yaml
```

```bash
module-b.yaml
└── module-d.yaml
```

As you can see, `module-d.yaml` is included twice, once via `module-a.yaml` and once via `module-b.yaml`.  Should `module-d.yaml` be executed twice? If so, in which order? Should it be executed only once? If so, which include should be used?

Currently, without recursive includes, we do not have these questions.

### 3. Cyclic dependencies

Although this behavior is often just an oversight, in the case of recursive includes, we open the door to cyclic dependencies. 

```bash
blueprint.yaml
└── module-a.yaml
```

```bash
module-a.yaml
└── module-b.yaml
```

```bash
module-b.yaml
└── module-a.yaml
```

As shown above, `module-a.yaml` includes `module-b.yaml`, which in turn includes `module-a.yaml`. This leeds to an infinite loop and will probably crash the engine - Something we currently completely avoid by not supporting recursive includes.

### Are recursive includes planned for the future?

Yes. As soon as projects start to grow and become more complex, nested or recursive includes will be a requirement. But as shown above, there are multiple pitfalls that need to be addressed before we can support this feature.