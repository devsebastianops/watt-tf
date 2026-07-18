# Terraform Modules

In modern cloud infrastructure, you rarely write every resource from scratch. Instead, platform teams rely on verified, reusable **Terraform Modules** (either from the public Terraform Registry or private corporate registries).

Watt TF integrates seamlessly with Terraform Modules. Since module calls are just JSON objects under the hood of `.tf.json`, your blueprints can dynamically stamp out module invocations based on your input data.


## 1. The Input Data (`input.yaml`)

Imagine your engineering teams need to provision static websites with Amazon S3 and CloudFront. They just define the core parameters of their sites:

```yaml
team: marketing
websites:
  - name: corporate-blog
    domain: blog.mycompany.com
    price_class: PriceClass_100 # US and Europe only (cost-efficient)
  - name: documentation
    domain: docs.mycompany.com
    price_class: PriceClass_All # Global distribution
```

## 2. The Blueprint (`blueprint.yaml`)

Instead of writing complex HCL wrapping logic, your blueprint uses `for_each` to iterate over the array of websites and invokes a verified CloudFront/S3 module for each entry.

```yaml
transform:
  # Iterate over all websites and generate a module block for each
  - for_each: input.websites
    target: "module.cloudfront_s3_website_${item.name}"
    value:
      # The standard source property Terraform expects
      source: "terraform-aws-modules/s3-bucket/aws//modules/notification" # Example registry path
      version: "4.1.0"

      # Passing dynamic values directly into the module variables
      bucket_name: "mycompany-web-${item.name}"
      domain_name: ${item.domain}
      price_class: ${item.price_class}

      tags:
        ManagedBy: "Watt TF"
        Owner: ${input.team}
```

## 3. Execution

When you run wtf build, Watt TF automatically expands the blueprint into the exact structural JSON format that Terraform needs to resolve and download the modules:

```bash
# 1. Generate the modules JSON
wtf build -i input.yaml -b blueprint.yaml -o static-sites.tf.json

# 2. Let Terraform fetch and apply the remote modules
terraform init
terraform apply -auto-approve
```

The generated `static-sites.tf.json` will look like this:

```json
{
  "module": {
    "cloudfront_s3_website_corporate_blog": {
      "source": "terraform-aws-modules/s3-bucket/aws//modules/notification",
      "version": "4.1.0",
      "bucket_name": "mycompany-web-corporate-blog",
      "domain_name": "blog.mycompany.com",
      "price_class": "PriceClass_100",
      "tags": {
        "ManagedBy": "Watt TF",
        "Owner": "marketing"
      }
    },
    "cloudfront_s3_website_documentation": {
      "source": "terraform-aws-modules/s3-bucket/aws//modules/notification",
      "version": "4.1.0",
      "bucket_name": "mycompany-web-documentation",
      "domain_name": "docs.mycompany.com",
      "price_class": "PriceClass_All",
      "tags": {
        "ManagedBy": "Watt TF",
        "Owner": "marketing"
      }
    }
  }
}
```