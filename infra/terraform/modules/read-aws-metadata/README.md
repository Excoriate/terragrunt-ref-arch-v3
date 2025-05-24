# Terraform Module: Read AWS Account Metadata

## Overview
> **Note:** This simple module uses AWS data sources (`aws_caller_identity`, `aws_partition`, `aws_region`) to retrieve metadata about the current AWS environment (Account ID, ARN, User ID, Partition, Region Name, Region Description) where Terraform is executing. It provides these details as outputs for use in other modules or configurations, often useful in Terragrunt setups.

### 🔑 Key Features
- Retrieves current AWS Account ID.
- Retrieves current AWS Caller ARN and User ID.
- Retrieves current AWS Partition (e.g., `aws`, `aws-us-gov`).
- Retrieves current AWS Region name and description.
- No resources are created, only data sources are used.

### 📋 Usage Guidelines
1. Include the module in your Terraform configuration.
2. Access the desired metadata via the module's outputs (e.g., `module.read_aws_account.account_id`).

```hcl
module "account_info" {
  source = "../path/to/read-aws-account"
  # No inputs required beyond provider configuration
}

resource "aws_s3_bucket" "example" {
  # Example using an output from this module
  bucket = "my-bucket-${module.account_info.account_id}-${module.account_info.region_name}"
}
```

<!-- BEGIN_TF_DOCS -->
<!-- END_TF_DOCS -->
