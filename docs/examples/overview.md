# Real world examples and How-tos

## Real world examples

This section contains real world examples of how to use the library in different scenarios. Each example is designed to demonstrate specific features and best practices for implementing the library effectively.

- **[Terraform Modules](./modules.md)** - Explore Watt TF using Terraform modules to use reusable infrastructure components.
- **[Multi Environment](./multi-environment.md)** - Showcase of Watt TF spinning up multiple environments through blueprints.
- **[Platform Engineering](./platform-engineering.md)** - Using Watt TF to build a platform engineering solution with reusable blueprints and modules.

## How-tos

This section contains more basic how-to's to demonstrate solutions for actual problems. Each one has a scenario and the solution with explanations, if needed.

- **[Dynamic targets](./dynamic-targets.md)** - Use interpolation to create dynamic target references.
- **[Deep Merging](./deep-merging.md)** - Explains the concept of extending already configured targets with more values via deep merging.
- **[Nested Targets](./nested-paths.md)** - How to define nested targets and span up complex objects.
- **[Preserve Types](./preserve-types.md)** - Demonstrates how types are preserved via CEL expressions.
- **[List Values](./list-values.md)** - How to work with list values inside blueprints.
- **[Environment Variables](./environment-variables.md)** - Using environment variables within Watt TF blueprints.
- **[Includes](./includes.md)** - How to include external blueprints within a blueprint.
- **[Escaping targets](./escaping-targets.md)** - Escaping target references that includes `.` characters.
- **[Terraform references](./terraform-references.md)** - Working with Terraform references like `var`, `module`, `local`, or `resource` within Watt TF blueprints.
- **[Loading external data](./loading-external-data.md)** - How to load and use external data within Watt TF blueprints.
- **[Reoccuring variables](./reoccuring-variables.md)** - Externalize complex and reoccuring variables for reuse across blueprints.
- **[Optional Values](./optional-values.md)** - Explains how to define and use optional values within Watt TF blueprints.
- **[Conditional Values](./conditional-values.md)** - Demonstrates how to define and use conditional values within Watt TF blueprints.
- **[Loops](./loops.md)** - How to define and use loops within Watt TF blueprints.
- **[Functions](./functions.md)** - How to use already provided CEL functions within Watt TF blueprints.
- **[Plugins](./plugins.md)** - How to use and manage plugins within Watt TF blueprints.
