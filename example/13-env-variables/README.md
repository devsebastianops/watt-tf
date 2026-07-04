# Example 13: Environment Variables

This example demonstrates how to use environment variables in CEL conditions and interpolations.

## What it does

- Uses the `env` CEL variable to access environment variables
- Evaluates a condition based on the `DEPLOYMENT_ENV` environment variable
- Interpolates values using environment variables (`DEPLOYMENT_ENV`, `DEPLOYMENT_REGION`)

## Environment Variables

This example requires the following environment variables to be set:
- `DEPLOYMENT_ENV=production`
- `DEPLOYMENT_REGION=us-east-1`

## Key Features

- **Condition with env variable:** `if: env.DEPLOYMENT_ENV == 'production'`
- **String interpolation:** `"${env.DEPLOYMENT_ENV}"` and `"${env.DEPLOYMENT_REGION}"`

## Expected Output

When `DEPLOYMENT_ENV=production` and `DEPLOYMENT_REGION=us-east-1`, the transformation creates a compute instance with the configuration values from both input and environment variables.
